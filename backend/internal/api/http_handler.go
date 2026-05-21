package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"ecochitas/internal/auth"
	"ecochitas/internal/domain"
	"ecochitas/internal/realtime"
	"ecochitas/internal/storage"
)

type Api_handler struct {
	bin_repository              *storage.Bin_repository
	truck_position_repository   *storage.Truck_position_repository
	route_repository            route_service
	recycling_zone_repository   *storage.Recycling_zone_repository
	operations_repository       operations_service
	truck_position_stream       *realtime.Truck_position_stream
	operations_event_publisher  operations_event_publisher_service
	jwt_authenticator           *auth.Jwt_authenticator
	auth_enable_dev_token_issue bool
	application_logger          *slog.Logger
}

func New_api_handler(
	bin_repository *storage.Bin_repository,
	truck_position_repository *storage.Truck_position_repository,
	route_repository route_service,
	recycling_zone_repository *storage.Recycling_zone_repository,
	operations_repository operations_service,
	truck_position_stream *realtime.Truck_position_stream,
	operations_event_publisher operations_event_publisher_service,
	jwt_authenticator *auth.Jwt_authenticator,
	auth_enable_dev_token_issue bool,
	application_logger *slog.Logger,
) *Api_handler {
	return &Api_handler{
		bin_repository:              bin_repository,
		truck_position_repository:   truck_position_repository,
		route_repository:            route_repository,
		recycling_zone_repository:   recycling_zone_repository,
		operations_repository:       operations_repository,
		truck_position_stream:       truck_position_stream,
		operations_event_publisher:  operations_event_publisher,
		jwt_authenticator:           jwt_authenticator,
		auth_enable_dev_token_issue: auth_enable_dev_token_issue,
		application_logger:          application_logger,
	}
}

func (api_handler *Api_handler) Register_routes(http_multiplexer *http.ServeMux) {
	http_multiplexer.HandleFunc("GET /healthz", api_handler.handle_health_check)
	http_multiplexer.HandleFunc("GET /v1/auth/me", api_handler.with_required_roles(
		[]string{"citizen", "driver", "admin", "condominium_admin"},
		api_handler.handle_get_authenticated_profile,
	))
	if api_handler.auth_enable_dev_token_issue {
		http_multiplexer.HandleFunc("POST /v1/auth/dev-token", api_handler.handle_issue_dev_access_token)
	}
	http_multiplexer.HandleFunc("GET /v1/bins", api_handler.handle_list_bins)
	http_multiplexer.HandleFunc("GET /v1/collection-routes", api_handler.handle_list_collection_routes)
	http_multiplexer.HandleFunc("GET /v1/collection-routes/{route_identifier}/stops", api_handler.handle_list_route_stops_by_route_identifier)
	http_multiplexer.HandleFunc("POST /v1/admin/collection-routes", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_create_collection_route,
	))
	http_multiplexer.HandleFunc("PATCH /v1/admin/collection-routes/{route_identifier}", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_update_collection_route,
	))
	http_multiplexer.HandleFunc("PUT /v1/admin/collection-routes/{route_identifier}/stops", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_sync_route_stops,
	))
	http_multiplexer.HandleFunc("GET /v1/admin/collection-routes/{route_identifier}/revisions", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_list_collection_route_revisions,
	))
	http_multiplexer.HandleFunc("POST /v1/admin/truck-route-assignments", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_create_truck_route_assignment,
	))
	http_multiplexer.HandleFunc("GET /v1/admin/truck-route-assignments", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_list_truck_route_assignments,
	))
	http_multiplexer.HandleFunc("GET /v1/admin/route-deviation-alerts", api_handler.with_required_roles(
		[]string{"admin"},
		api_handler.handle_list_route_deviation_alerts,
	))
	http_multiplexer.HandleFunc("GET /v1/trucks/latest-position", api_handler.handle_get_latest_truck_position)
	http_multiplexer.HandleFunc("GET /v1/trucks/latest-positions", api_handler.handle_list_latest_truck_positions)
	http_multiplexer.HandleFunc("GET /v1/trucks/stream", api_handler.handle_stream_truck_positions)
	http_multiplexer.HandleFunc("GET /v1/driver/route-deviation", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_get_truck_route_deviation,
	))
	http_multiplexer.HandleFunc("POST /v1/driver/route-deviation-alerts", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_create_route_deviation_alert,
	))
	http_multiplexer.HandleFunc("POST /v1/bins/sensor-events", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_ingest_bin_sensor_event,
	))
	http_multiplexer.HandleFunc("POST /v1/driver/collection-events", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_record_driver_collection_event,
	))
	http_multiplexer.HandleFunc("POST /v1/driver/route-blockages", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_create_route_blockage_report,
	))
	http_multiplexer.HandleFunc("GET /v1/driver/route-blockages", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_list_route_blockage_reports,
	))
	http_multiplexer.HandleFunc("PATCH /v1/driver/route-blockages/{blockage_identifier}", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_update_route_blockage_report_status,
	))
	http_multiplexer.HandleFunc("GET /v1/recycling/containers", api_handler.handle_list_zone_recycling_containers)
	http_multiplexer.HandleFunc("POST /v1/recycling/cycles/start", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_start_recycling_collection_cycle,
	))
	http_multiplexer.HandleFunc("POST /v1/recycling/evidence-submissions", api_handler.with_required_roles(
		[]string{"citizen", "driver", "admin", "condominium_admin"},
		api_handler.handle_submit_recycling_evidence,
	))
	http_multiplexer.HandleFunc("POST /v1/recycling/cycles/close", api_handler.with_required_roles(
		[]string{"driver", "admin"},
		api_handler.handle_close_recycling_collection_cycle,
	))
	http_multiplexer.HandleFunc("GET /v1/recycling/cycles/summary", api_handler.handle_get_recycling_cycle_summary)
}

func (api_handler *Api_handler) handle_health_check(response_writer http.ResponseWriter, request *http.Request) {
	health_response_payload := map[string]any{
		"status":    "ok",
		"service":   "ecochitas_backend",
		"timestamp": time.Now().UTC(),
	}
	write_json_response(response_writer, http.StatusOK, health_response_payload)
}

func (api_handler *Api_handler) handle_get_authenticated_profile(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"user_identifier": auth_claims.User_identifier,
		"role_name":       auth_claims.Role_name,
		"full_name":       auth_claims.Full_name,
		"subject":         auth_claims.Subject,
		"issuer":          auth_claims.Issuer,
		"audience":        auth_claims.Audience,
		"expires_at":      auth_claims.ExpiresAt,
	})
}

func (api_handler *Api_handler) handle_issue_dev_access_token(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if !api_handler.auth_enable_dev_token_issue {
		write_json_error(response_writer, http.StatusNotFound, "resource_not_found")
		return
	}

	var issue_access_token_payload struct {
		User_identifier string `json:"user_identifier"`
		Role_name       string `json:"role_name"`
		Full_name       string `json:"full_name"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&issue_access_token_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_issue_dev_access_token_payload")
		return
	}

	generate_access_token_command := auth.Generate_access_token_command{
		User_identifier: strings.TrimSpace(issue_access_token_payload.User_identifier),
		Role_name:       strings.TrimSpace(issue_access_token_payload.Role_name),
		Full_name:       strings.TrimSpace(issue_access_token_payload.Full_name),
	}
	access_token, expires_at, generate_access_token_error := api_handler.jwt_authenticator.Generate_access_token(
		generate_access_token_command,
	)
	if generate_access_token_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, generate_access_token_error.Error())
		return
	}

	write_json_response(response_writer, http.StatusCreated, map[string]any{
		"access_token": access_token,
		"token_type":   "Bearer",
		"expires_at":   expires_at,
		"role_name":    generate_access_token_command.Role_name,
	})
}

func (api_handler *Api_handler) handle_list_bins(response_writer http.ResponseWriter, request *http.Request) {
	query_context := request.Context()

	bin_status_list, load_bins_error := api_handler.bin_repository.List_bins(query_context)
	if load_bins_error != nil {
		api_handler.application_logger.Error("failed_to_load_bins", "error", load_bins_error)
		write_json_error(response_writer, http.StatusInternalServerError, "failed_to_load_bins")
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"items": bin_status_list,
		"total": len(bin_status_list),
	})
}

func (api_handler *Api_handler) handle_get_latest_truck_position(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	raw_truck_identifier := request.URL.Query().Get("truck_identifier")
	truck_identifier := strings.TrimSpace(raw_truck_identifier)

	if truck_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "truck_identifier_query_param_is_required")
		return
	}

	query_context := request.Context()
	latest_truck_position, query_position_error := api_handler.truck_position_repository.Get_latest_truck_position(
		query_context,
		truck_identifier,
	)
	if query_position_error != nil {
		api_handler.application_logger.Error(
			"failed_to_load_latest_truck_position",
			"error",
			query_position_error,
			"truck_identifier",
			truck_identifier,
		)
		write_json_error(response_writer, http.StatusInternalServerError, "failed_to_load_latest_truck_position")
		return
	}

	if latest_truck_position == nil {
		write_json_error(response_writer, http.StatusNotFound, "truck_position_not_found")
		return
	}

	write_json_response(response_writer, http.StatusOK, latest_truck_position)
}

func (api_handler *Api_handler) handle_list_latest_truck_positions(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	latest_truck_positions, load_positions_error := api_handler.truck_position_repository.List_latest_truck_positions(
		query_context,
	)
	if load_positions_error != nil {
		api_handler.application_logger.Error("failed_to_load_latest_truck_positions", "error", load_positions_error)
		write_json_error(response_writer, http.StatusInternalServerError, "failed_to_load_latest_truck_positions")
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"items": latest_truck_positions,
		"total": len(latest_truck_positions),
	})
}

func (api_handler *Api_handler) handle_stream_truck_positions(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	response_flusher, supports_flusher := response_writer.(http.Flusher)
	if !supports_flusher {
		write_json_error(response_writer, http.StatusInternalServerError, "streaming_is_not_supported")
		return
	}

	raw_truck_identifier_filter := request.URL.Query().Get("truck_identifier")
	truck_identifier_filter := strings.TrimSpace(raw_truck_identifier_filter)

	response_writer.Header().Set("Content-Type", "text/event-stream")
	response_writer.Header().Set("Cache-Control", "no-cache")
	response_writer.Header().Set("Connection", "keep-alive")
	response_writer.Header().Set("X-Accel-Buffering", "no")

	subscriber_identifier, truck_position_event_ch, unsubscribe_function := api_handler.truck_position_stream.Subscribe(
		truck_identifier_filter,
	)
	defer unsubscribe_function()

	ready_event_payload := map[string]any{
		"subscriber_identifier":   subscriber_identifier,
		"truck_identifier_filter": truck_identifier_filter,
		"connected_at":            time.Now().UTC(),
	}
	write_ready_event_error := write_sse_event(response_writer, response_flusher, "ready", ready_event_payload)
	if write_ready_event_error != nil {
		api_handler.application_logger.Warn(
			"truck_position_stream_ready_event_write_failed",
			"error",
			write_ready_event_error,
			"subscriber_identifier",
			subscriber_identifier,
		)
		return
	}

	keep_alive_ticker := time.NewTicker(15 * time.Second)
	defer keep_alive_ticker.Stop()

	request_context := request.Context()
	for {
		select {
		case <-request_context.Done():
			return
		case <-keep_alive_ticker.C:
			write_sse_comment(response_writer, response_flusher, "keep_alive")
		case truck_latest_position, channel_open := <-truck_position_event_ch:
			if !channel_open {
				return
			}

			write_stream_error := write_sse_event(response_writer, response_flusher, "truck_position", truck_latest_position)
			if write_stream_error != nil {
				api_handler.application_logger.Warn(
					"truck_position_stream_write_failed",
					"error",
					write_stream_error,
					"subscriber_identifier",
					subscriber_identifier,
				)
				return
			}
		}
	}
}

func write_sse_event(
	response_writer http.ResponseWriter,
	response_flusher http.Flusher,
	event_name string,
	event_payload any,
) error {
	serialized_payload, marshal_payload_error := json.Marshal(event_payload)
	if marshal_payload_error != nil {
		return fmt.Errorf("failed_to_marshal_sse_payload: %w", marshal_payload_error)
	}

	_, write_event_name_error := response_writer.Write([]byte("event: " + event_name + "\n"))
	if write_event_name_error != nil {
		return write_event_name_error
	}

	_, write_event_payload_error := response_writer.Write([]byte("data: " + string(serialized_payload) + "\n\n"))
	if write_event_payload_error != nil {
		return write_event_payload_error
	}

	response_flusher.Flush()
	return nil
}

func write_sse_comment(
	response_writer http.ResponseWriter,
	response_flusher http.Flusher,
	comment_text string,
) {
	_, _ = response_writer.Write([]byte(": " + comment_text + "\n\n"))
	response_flusher.Flush()
}

func (api_handler *Api_handler) handle_list_zone_recycling_containers(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()
	zone_code_filter := strings.TrimSpace(request.URL.Query().Get("zone_code"))

	container_view_list, list_containers_error := api_handler.recycling_zone_repository.List_zone_recycling_containers(
		query_context,
		zone_code_filter,
	)
	if list_containers_error != nil {
		api_handler.application_logger.Error("failed_to_list_zone_recycling_containers", "error", list_containers_error)
		write_json_error(response_writer, http.StatusInternalServerError, "failed_to_list_zone_recycling_containers")
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"items": container_view_list,
		"total": len(container_view_list),
	})
}

func (api_handler *Api_handler) handle_start_recycling_collection_cycle(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	var recycling_cycle_start_command_payload struct {
		Container_identifier      string `json:"container_identifier"`
		Scheduled_collection_date string `json:"scheduled_collection_date"`
		Collection_operator_name  string `json:"collection_operator_name"`
	}
	decode_request_error := json.NewDecoder(request.Body).Decode(&recycling_cycle_start_command_payload)
	if decode_request_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_recycling_cycle_start_payload")
		return
	}

	container_identifier := strings.TrimSpace(recycling_cycle_start_command_payload.Container_identifier)
	if container_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "container_identifier_is_required")
		return
	}

	var scheduled_collection_date time.Time
	if strings.TrimSpace(recycling_cycle_start_command_payload.Scheduled_collection_date) == "" {
		scheduled_collection_date = time.Now().UTC()
	} else {
		parsed_collection_date, parse_collection_date_error := time.Parse(
			"2006-01-02",
			strings.TrimSpace(recycling_cycle_start_command_payload.Scheduled_collection_date),
		)
		if parse_collection_date_error != nil {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_scheduled_collection_date_format")
			return
		}
		scheduled_collection_date = parsed_collection_date.UTC()
	}

	recycling_cycle_start_command := domain.Recycling_cycle_start_command{
		Container_identifier:      container_identifier,
		Scheduled_collection_date: scheduled_collection_date,
		Collection_operator_name:  strings.TrimSpace(recycling_cycle_start_command_payload.Collection_operator_name),
	}

	created_recycling_cycle, start_cycle_error := api_handler.recycling_zone_repository.Start_recycling_collection_cycle(
		query_context,
		recycling_cycle_start_command,
	)
	if start_cycle_error != nil {
		api_handler.application_logger.Error("failed_to_start_recycling_collection_cycle", "error", start_cycle_error)
		write_json_error(response_writer, http.StatusBadRequest, "failed_to_start_recycling_collection_cycle")
		return
	}

	write_json_response(response_writer, http.StatusCreated, created_recycling_cycle)
}

func (api_handler *Api_handler) handle_submit_recycling_evidence(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	var recycling_evidence_submission_payload struct {
		Cycle_identifier     string  `json:"cycle_identifier"`
		Household_identifier string  `json:"household_identifier"`
		Evidence_photo_url   string  `json:"evidence_photo_url"`
		Evidence_captured_at string  `json:"evidence_captured_at"`
		Evidence_latitude    float64 `json:"evidence_latitude"`
		Evidence_longitude   float64 `json:"evidence_longitude"`
		Validation_status    string  `json:"validation_status"`
		Rejection_reason     string  `json:"rejection_reason"`
	}
	decode_request_error := json.NewDecoder(request.Body).Decode(&recycling_evidence_submission_payload)
	if decode_request_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_recycling_evidence_submission_payload")
		return
	}

	if strings.TrimSpace(recycling_evidence_submission_payload.Cycle_identifier) == "" {
		write_json_error(response_writer, http.StatusBadRequest, "cycle_identifier_is_required")
		return
	}

	if strings.TrimSpace(recycling_evidence_submission_payload.Household_identifier) == "" {
		write_json_error(response_writer, http.StatusBadRequest, "household_identifier_is_required")
		return
	}

	if strings.TrimSpace(recycling_evidence_submission_payload.Evidence_photo_url) == "" {
		write_json_error(response_writer, http.StatusBadRequest, "evidence_photo_url_is_required")
		return
	}

	if recycling_evidence_submission_payload.Evidence_latitude < -90 || recycling_evidence_submission_payload.Evidence_latitude > 90 {
		write_json_error(response_writer, http.StatusBadRequest, "evidence_latitude_out_of_range")
		return
	}
	if recycling_evidence_submission_payload.Evidence_longitude < -180 || recycling_evidence_submission_payload.Evidence_longitude > 180 {
		write_json_error(response_writer, http.StatusBadRequest, "evidence_longitude_out_of_range")
		return
	}

	var evidence_captured_at time.Time
	if strings.TrimSpace(recycling_evidence_submission_payload.Evidence_captured_at) == "" {
		evidence_captured_at = time.Now().UTC()
	} else {
		parsed_evidence_captured_at, parse_evidence_captured_at_error := time.Parse(
			time.RFC3339,
			strings.TrimSpace(recycling_evidence_submission_payload.Evidence_captured_at),
		)
		if parse_evidence_captured_at_error != nil {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_evidence_captured_at_format")
			return
		}
		evidence_captured_at = parsed_evidence_captured_at.UTC()
	}

	recycling_evidence_submission_command := domain.Recycling_evidence_submission_command{
		Cycle_identifier:     strings.TrimSpace(recycling_evidence_submission_payload.Cycle_identifier),
		Household_identifier: strings.TrimSpace(recycling_evidence_submission_payload.Household_identifier),
		Evidence_photo_url:   strings.TrimSpace(recycling_evidence_submission_payload.Evidence_photo_url),
		Evidence_captured_at: evidence_captured_at,
		Evidence_latitude:    recycling_evidence_submission_payload.Evidence_latitude,
		Evidence_longitude:   recycling_evidence_submission_payload.Evidence_longitude,
		Validation_status:    strings.TrimSpace(recycling_evidence_submission_payload.Validation_status),
		Rejection_reason:     strings.TrimSpace(recycling_evidence_submission_payload.Rejection_reason),
	}

	created_or_updated_submission, submit_recycling_evidence_error := api_handler.recycling_zone_repository.Submit_recycling_evidence(
		query_context,
		recycling_evidence_submission_command,
	)
	if submit_recycling_evidence_error != nil {
		api_handler.application_logger.Error("failed_to_submit_recycling_evidence", "error", submit_recycling_evidence_error)
		write_json_error(response_writer, http.StatusBadRequest, submit_recycling_evidence_error.Error())
		return
	}

	write_json_response(response_writer, http.StatusCreated, created_or_updated_submission)
}

func (api_handler *Api_handler) handle_close_recycling_collection_cycle(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	var recycling_cycle_close_payload struct {
		Cycle_identifier                  string   `json:"cycle_identifier"`
		Raw_points_total                  float64  `json:"raw_points_total"`
		Contamination_level               string   `json:"contamination_level"`
		Contamination_discount_percentage *float64 `json:"contamination_discount_percentage"`
		Contamination_notes               string   `json:"contamination_notes"`
		Collection_operator_name          string   `json:"collection_operator_name"`
	}
	decode_request_error := json.NewDecoder(request.Body).Decode(&recycling_cycle_close_payload)
	if decode_request_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_recycling_cycle_close_payload")
		return
	}

	recycling_cycle_close_command := domain.Recycling_cycle_close_command{
		Cycle_identifier:                  strings.TrimSpace(recycling_cycle_close_payload.Cycle_identifier),
		Raw_points_total:                  recycling_cycle_close_payload.Raw_points_total,
		Contamination_level:               strings.TrimSpace(recycling_cycle_close_payload.Contamination_level),
		Contamination_discount_percentage: recycling_cycle_close_payload.Contamination_discount_percentage,
		Contamination_notes:               strings.TrimSpace(recycling_cycle_close_payload.Contamination_notes),
		Collection_operator_name:          strings.TrimSpace(recycling_cycle_close_payload.Collection_operator_name),
	}

	closed_cycle_summary, close_cycle_error := api_handler.recycling_zone_repository.Close_recycling_collection_cycle(
		query_context,
		recycling_cycle_close_command,
	)
	if close_cycle_error != nil {
		api_handler.application_logger.Error("failed_to_close_recycling_collection_cycle", "error", close_cycle_error)
		write_json_error(response_writer, http.StatusBadRequest, close_cycle_error.Error())
		return
	}

	write_json_response(response_writer, http.StatusOK, closed_cycle_summary)
}

func (api_handler *Api_handler) handle_get_recycling_cycle_summary(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()
	cycle_identifier := strings.TrimSpace(request.URL.Query().Get("cycle_identifier"))
	if cycle_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "cycle_identifier_query_param_is_required")
		return
	}

	recycling_cycle_summary, load_summary_error := api_handler.recycling_zone_repository.Get_recycling_cycle_summary(
		query_context,
		cycle_identifier,
	)
	if load_summary_error != nil {
		api_handler.application_logger.Error("failed_to_load_recycling_cycle_summary", "error", load_summary_error)
		write_json_error(response_writer, http.StatusBadRequest, load_summary_error.Error())
		return
	}

	write_json_response(response_writer, http.StatusOK, recycling_cycle_summary)
}
