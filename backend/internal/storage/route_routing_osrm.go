package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"ecochitas/internal/domain"
)

const (
	default_osrm_base_url             = "https://router.project-osrm.org"
	default_osrm_request_timeout      = 8 * time.Second
	osrm_minimum_request_interval     = 1 * time.Second
	maximum_osrm_error_body_byte_size = 1024
)

type osrm_route_service struct {
	base_url             string
	http_client          *http.Client
	request_rate_limiter *osrm_request_rate_limiter
}

type osrm_request_rate_limiter struct {
	mutex                         sync.Mutex
	next_allowed_request_time_utc time.Time
}

type osrm_route_response_payload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Routes  []struct {
		Geometry struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func new_osrm_route_service(raw_osrm_base_url string) *osrm_route_service {
	normalized_osrm_base_url := normalize_osrm_base_url(raw_osrm_base_url)

	return &osrm_route_service{
		base_url: normalized_osrm_base_url,
		http_client: &http.Client{
			Timeout: default_osrm_request_timeout,
		},
		request_rate_limiter: new_osrm_request_rate_limiter(),
	}
}

func new_osrm_request_rate_limiter() *osrm_request_rate_limiter {
	return &osrm_request_rate_limiter{}
}

func (osrm_service *osrm_route_service) calculate_road_path_coordinates(
	application_context context.Context,
	ordered_route_stop_coordinate_list []domain.Route_path_coordinate,
) ([]domain.Route_road_path_coordinate, error) {
	if len(ordered_route_stop_coordinate_list) < 2 {
		return nil, fmt.Errorf("insufficient_route_stops_for_osrm")
	}

	wait_for_request_slot_error := osrm_service.request_rate_limiter.wait_for_request_slot(application_context)
	if wait_for_request_slot_error != nil {
		return nil, fmt.Errorf("failed_to_wait_for_osrm_request_slot: %w", wait_for_request_slot_error)
	}

	coordinate_segment_builder := strings.Builder{}
	for stop_index, route_stop_coordinate := range ordered_route_stop_coordinate_list {
		if stop_index > 0 {
			coordinate_segment_builder.WriteString(";")
		}
		coordinate_segment_builder.WriteString(
			strconv.FormatFloat(route_stop_coordinate.Longitude, 'f', 7, 64),
		)
		coordinate_segment_builder.WriteString(",")
		coordinate_segment_builder.WriteString(
			strconv.FormatFloat(route_stop_coordinate.Latitude, 'f', 7, 64),
		)
	}

	request_endpoint_url := osrm_service.base_url +
		"/route/v1/driving/" +
		coordinate_segment_builder.String() +
		"?overview=full&geometries=geojson"
	request_with_context, create_request_error := http.NewRequestWithContext(
		application_context,
		http.MethodGet,
		request_endpoint_url,
		nil,
	)
	if create_request_error != nil {
		return nil, fmt.Errorf("failed_to_create_osrm_request: %w", create_request_error)
	}

	osrm_response, request_error := osrm_service.http_client.Do(request_with_context)
	if request_error != nil {
		return nil, fmt.Errorf("failed_to_execute_osrm_request: %w", request_error)
	}
	defer osrm_response.Body.Close()

	if osrm_response.StatusCode != http.StatusOK {
		limited_error_body, _ := io.ReadAll(io.LimitReader(osrm_response.Body, maximum_osrm_error_body_byte_size))
		return nil, fmt.Errorf(
			"osrm_non_success_status status=%d body=%s",
			osrm_response.StatusCode,
			strings.TrimSpace(string(limited_error_body)),
		)
	}

	response_payload := osrm_route_response_payload{}
	decode_payload_error := json.NewDecoder(osrm_response.Body).Decode(&response_payload)
	if decode_payload_error != nil {
		return nil, fmt.Errorf("failed_to_decode_osrm_response: %w", decode_payload_error)
	}

	if response_payload.Code != "Ok" {
		return nil, fmt.Errorf(
			"osrm_response_code_not_ok code=%s message=%s",
			strings.TrimSpace(response_payload.Code),
			strings.TrimSpace(response_payload.Message),
		)
	}
	if len(response_payload.Routes) == 0 {
		return nil, fmt.Errorf("osrm_route_response_has_no_routes")
	}

	route_geometry := response_payload.Routes[0].Geometry
	if route_geometry.Type != "LineString" {
		return nil, fmt.Errorf("osrm_geometry_type_not_supported")
	}
	if len(route_geometry.Coordinates) < 2 {
		return nil, fmt.Errorf("osrm_route_geometry_has_insufficient_coordinates")
	}

	road_path_coordinate_list := make(
		[]domain.Route_road_path_coordinate,
		0,
		len(route_geometry.Coordinates),
	)
	for _, raw_coordinate_pair := range route_geometry.Coordinates {
		if len(raw_coordinate_pair) < 2 {
			return nil, fmt.Errorf("osrm_coordinate_pair_is_invalid")
		}
		road_path_coordinate_list = append(
			road_path_coordinate_list,
			domain.Route_road_path_coordinate{
				raw_coordinate_pair[0],
				raw_coordinate_pair[1],
			},
		)
	}

	return road_path_coordinate_list, nil
}

func (request_rate_limiter *osrm_request_rate_limiter) wait_for_request_slot(
	application_context context.Context,
) error {
	request_rate_limiter.mutex.Lock()
	current_time_utc := time.Now().UTC()
	wait_duration := time.Duration(0)
	reserved_request_time_utc := current_time_utc
	if current_time_utc.Before(request_rate_limiter.next_allowed_request_time_utc) {
		wait_duration = request_rate_limiter.next_allowed_request_time_utc.Sub(current_time_utc)
		reserved_request_time_utc = request_rate_limiter.next_allowed_request_time_utc
	}
	request_rate_limiter.next_allowed_request_time_utc = reserved_request_time_utc.Add(
		osrm_minimum_request_interval,
	)
	request_rate_limiter.mutex.Unlock()

	if wait_duration <= 0 {
		return nil
	}

	wait_timer := time.NewTimer(wait_duration)
	defer wait_timer.Stop()

	select {
	case <-application_context.Done():
		return application_context.Err()
	case <-wait_timer.C:
		return nil
	}
}

func normalize_osrm_base_url(raw_osrm_base_url string) string {
	normalized_osrm_base_url := strings.TrimSpace(raw_osrm_base_url)
	if normalized_osrm_base_url == "" {
		normalized_osrm_base_url = default_osrm_base_url
	}

	return strings.TrimRight(normalized_osrm_base_url, "/")
}
