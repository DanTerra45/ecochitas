package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ecochitas/internal/domain"
)

func (api_handler *Api_handler) handle_list_collection_routes(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	query_context := request.Context()
	collection_weekday_filter, weekday_error := parse_optional_collection_weekday_query_parameter(
		strings.TrimSpace(request.URL.Query().Get("collection_weekday")),
	)
	if weekday_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_collection_weekday_query_param")
		return
	}

	has_is_active_filter, is_active_filter_value, is_active_filter_error := parse_optional_bool_query_parameter(
		strings.TrimSpace(request.URL.Query().Get("is_active")),
	)
	if is_active_filter_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_is_active_query_param")
		return
	}

	collection_route_list_query := domain.Collection_route_list_query{
		Zone_name_filter:          strings.TrimSpace(request.URL.Query().Get("zone_name")),
		Collection_weekday_filter: collection_weekday_filter,
		Has_is_active_filter:      has_is_active_filter,
		Is_active_filter_value:    is_active_filter_value,
	}

	collection_route_view_list, list_routes_error := api_handler.route_repository.List_collection_routes(
		query_context,
		collection_route_list_query,
	)
	if list_routes_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(list_routes_error)
		api_handler.application_logger.Error("failed_to_list_collection_routes", "error", list_routes_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"items": collection_route_view_list,
		"total": len(collection_route_view_list),
	})
}

func (api_handler *Api_handler) handle_create_collection_route(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	var create_collection_route_payload struct {
		Route_code         string `json:"route_code"`
		Route_name         string `json:"route_name"`
		Zone_name          string `json:"zone_name"`
		Collection_weekday int    `json:"collection_weekday"`
		Is_active          *bool  `json:"is_active"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&create_collection_route_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_create_collection_route_payload")
		return
	}

	is_active := true
	if create_collection_route_payload.Is_active != nil {
		is_active = *create_collection_route_payload.Is_active
	}

	collection_route_create_command := domain.Collection_route_create_command{
		Route_code:                    create_collection_route_payload.Route_code,
		Route_name:                    create_collection_route_payload.Route_name,
		Zone_name:                     create_collection_route_payload.Zone_name,
		Collection_weekday:            create_collection_route_payload.Collection_weekday,
		Is_active:                     is_active,
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
	}

	created_collection_route, create_collection_route_error := api_handler.route_repository.Create_collection_route(
		request.Context(),
		collection_route_create_command,
	)
	if create_collection_route_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(create_collection_route_error)
		api_handler.application_logger.Error(
			"failed_to_create_collection_route",
			"error",
			create_collection_route_error,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusCreated, created_collection_route)
}

func (api_handler *Api_handler) handle_update_collection_route(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	route_identifier := strings.TrimSpace(request.PathValue("route_identifier"))
	if route_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "route_identifier_path_param_is_required")
		return
	}

	var update_collection_route_payload struct {
		Route_code         *string `json:"route_code"`
		Route_name         *string `json:"route_name"`
		Zone_name          *string `json:"zone_name"`
		Collection_weekday *int    `json:"collection_weekday"`
		Is_active          *bool   `json:"is_active"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&update_collection_route_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_update_collection_route_payload")
		return
	}

	collection_route_update_command := domain.Collection_route_update_command{
		Route_identifier:              route_identifier,
		Route_code:                    update_collection_route_payload.Route_code,
		Route_name:                    update_collection_route_payload.Route_name,
		Zone_name:                     update_collection_route_payload.Zone_name,
		Collection_weekday:            update_collection_route_payload.Collection_weekday,
		Is_active:                     update_collection_route_payload.Is_active,
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
	}

	updated_collection_route, update_collection_route_error := api_handler.route_repository.Update_collection_route(
		request.Context(),
		collection_route_update_command,
	)
	if update_collection_route_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(update_collection_route_error)
		api_handler.application_logger.Error(
			"failed_to_update_collection_route",
			"error",
			update_collection_route_error,
			"route_identifier",
			route_identifier,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, updated_collection_route)
}

func (api_handler *Api_handler) handle_list_route_stops_by_route_identifier(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	query_context := request.Context()
	route_identifier := strings.TrimSpace(request.PathValue("route_identifier"))
	if route_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "route_identifier_path_param_is_required")
		return
	}

	route_stop_view_list, list_stops_error := api_handler.route_repository.List_route_stops_by_route_identifier(
		query_context,
		route_identifier,
	)
	if list_stops_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(list_stops_error)
		api_handler.application_logger.Error(
			"failed_to_list_route_stops_by_route_identifier",
			"error",
			list_stops_error,
			"route_identifier",
			route_identifier,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"route_identifier": route_identifier,
		"items":            route_stop_view_list,
		"total":            len(route_stop_view_list),
	})
}

func (api_handler *Api_handler) handle_sync_route_stops(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	route_identifier := strings.TrimSpace(request.PathValue("route_identifier"))
	if route_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "route_identifier_path_param_is_required")
		return
	}

	var sync_route_stops_payload struct {
		Stop_item_list []domain.Route_stop_sync_item `json:"stops"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&sync_route_stops_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_sync_route_stops_payload")
		return
	}

	route_stop_sync_command := domain.Route_stop_sync_command{
		Route_identifier:              route_identifier,
		Stop_item_list:                sync_route_stops_payload.Stop_item_list,
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
	}

	updated_route_stop_view_list, sync_route_stops_error := api_handler.route_repository.Sync_route_stops(
		request.Context(),
		route_stop_sync_command,
	)
	if sync_route_stops_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(sync_route_stops_error)
		api_handler.application_logger.Error(
			"failed_to_sync_route_stops",
			"error",
			sync_route_stops_error,
			"route_identifier",
			route_identifier,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, map[string]any{
		"route_identifier": route_identifier,
		"items":            updated_route_stop_view_list,
		"total":            len(updated_route_stop_view_list),
	})
}

func (api_handler *Api_handler) handle_list_collection_route_revisions(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	route_identifier := strings.TrimSpace(request.PathValue("route_identifier"))
	if route_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "route_identifier_path_param_is_required")
		return
	}

	limit_query_value := strings.TrimSpace(request.URL.Query().Get("limit"))
	limit := 0
	if limit_query_value != "" {
		parsed_limit, parse_limit_error := strconv.Atoi(limit_query_value)
		if parse_limit_error != nil || parsed_limit <= 0 {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_limit_query_param")
			return
		}
		limit = parsed_limit
	}

	collection_route_revision_list_result, list_route_revisions_error := api_handler.route_repository.List_collection_route_revisions(
		request.Context(),
		domain.Collection_route_revision_list_query{
			Route_identifier: route_identifier,
			Limit:            limit,
		},
	)
	if list_route_revisions_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(list_route_revisions_error)
		api_handler.application_logger.Error(
			"failed_to_list_collection_route_revisions",
			"error",
			list_route_revisions_error,
			"route_identifier",
			route_identifier,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, collection_route_revision_list_result)
}

func (api_handler *Api_handler) handle_get_truck_route_deviation(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	truck_identifier := strings.TrimSpace(request.URL.Query().Get("truck_identifier"))
	route_identifier := strings.TrimSpace(request.URL.Query().Get("route_identifier"))
	if truck_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "truck_identifier_query_param_is_required")
		return
	}

	deviation_threshold_meters := 0.0
	raw_deviation_threshold := strings.TrimSpace(request.URL.Query().Get("deviation_threshold_meters"))
	if raw_deviation_threshold != "" {
		parsed_deviation_threshold, parse_deviation_threshold_error := strconv.ParseFloat(raw_deviation_threshold, 64)
		if parse_deviation_threshold_error != nil || parsed_deviation_threshold <= 0 {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_deviation_threshold_meters_query_param")
			return
		}
		deviation_threshold_meters = parsed_deviation_threshold
	}

	truck_route_deviation_view, get_truck_route_deviation_error := api_handler.route_repository.Get_truck_route_deviation(
		request.Context(),
		domain.Truck_route_deviation_query{
			Truck_identifier:           truck_identifier,
			Route_identifier:           route_identifier,
			Deviation_threshold_meters: deviation_threshold_meters,
		},
	)
	if get_truck_route_deviation_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(get_truck_route_deviation_error)
		api_handler.application_logger.Error(
			"failed_to_get_truck_route_deviation",
			"error",
			get_truck_route_deviation_error,
			"truck_identifier",
			truck_identifier,
			"route_identifier",
			route_identifier,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, truck_route_deviation_view)
}

func (api_handler *Api_handler) handle_create_truck_route_assignment(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	var create_assignment_payload struct {
		Truck_identifier string `json:"truck_identifier"`
		Route_identifier string `json:"route_identifier"`
		Assignment_notes string `json:"assignment_notes"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&create_assignment_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_create_truck_route_assignment_payload")
		return
	}

	created_assignment, create_assignment_error := api_handler.route_repository.Create_truck_route_assignment(
		request.Context(),
		domain.Truck_route_assignment_create_command{
			Truck_identifier:            create_assignment_payload.Truck_identifier,
			Route_identifier:            create_assignment_payload.Route_identifier,
			Assigned_by_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
			Assignment_notes:            create_assignment_payload.Assignment_notes,
		},
	)
	if create_assignment_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(create_assignment_error)
		api_handler.application_logger.Error("failed_to_create_truck_route_assignment", "error", create_assignment_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusCreated, created_assignment)
}

func (api_handler *Api_handler) handle_list_truck_route_assignments(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	var has_is_active_filter bool
	var is_active_filter_value bool
	raw_is_active_filter_value := strings.TrimSpace(request.URL.Query().Get("is_active"))
	if raw_is_active_filter_value != "" {
		parsed_is_active_filter_value, parse_is_active_filter_error := strconv.ParseBool(raw_is_active_filter_value)
		if parse_is_active_filter_error != nil {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_is_active_query_param")
			return
		}
		has_is_active_filter = true
		is_active_filter_value = parsed_is_active_filter_value
	}

	limit_value := 0
	raw_limit_value := strings.TrimSpace(request.URL.Query().Get("limit"))
	if raw_limit_value != "" {
		parsed_limit_value, parse_limit_error := strconv.Atoi(raw_limit_value)
		if parse_limit_error != nil || parsed_limit_value <= 0 {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_limit_query_param")
			return
		}
		limit_value = parsed_limit_value
	}

	var is_active_filter *bool
	if has_is_active_filter {
		is_active_filter = &is_active_filter_value
	}

	truck_route_assignment_list_result, list_assignments_error := api_handler.route_repository.List_truck_route_assignments(
		request.Context(),
		domain.Truck_route_assignment_list_query{
			Truck_identifier: strings.TrimSpace(request.URL.Query().Get("truck_identifier")),
			Route_identifier: strings.TrimSpace(request.URL.Query().Get("route_identifier")),
			Is_active_filter: is_active_filter,
			Limit:            limit_value,
		},
	)
	if list_assignments_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(list_assignments_error)
		api_handler.application_logger.Error("failed_to_list_truck_route_assignments", "error", list_assignments_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, truck_route_assignment_list_result)
}

func (api_handler *Api_handler) handle_create_route_deviation_alert(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	var create_alert_payload struct {
		Truck_identifier           string  `json:"truck_identifier"`
		Route_identifier           string  `json:"route_identifier"`
		Deviation_threshold_meters float64 `json:"deviation_threshold_meters"`
		Alert_notes                string  `json:"alert_notes"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&create_alert_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_create_route_deviation_alert_payload")
		return
	}

	created_route_deviation_alert_record, create_alert_error := api_handler.route_repository.Create_route_deviation_alert(
		request.Context(),
		domain.Route_deviation_alert_create_command{
			Truck_identifier:             create_alert_payload.Truck_identifier,
			Route_identifier:             create_alert_payload.Route_identifier,
			Deviation_threshold_meters:   create_alert_payload.Deviation_threshold_meters,
			Triggered_by_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
			Alert_notes:                  create_alert_payload.Alert_notes,
		},
	)
	if create_alert_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(create_alert_error)
		api_handler.application_logger.Error("failed_to_create_route_deviation_alert", "error", create_alert_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	if api_handler.operations_event_publisher != nil {
		publish_event_error := api_handler.operations_event_publisher.Publish_route_deviation_alert_event(
			created_route_deviation_alert_record,
		)
		if publish_event_error != nil {
			api_handler.application_logger.Warn("failed_to_publish_route_deviation_alert_event", "error", publish_event_error)
		}
	}

	write_json_response(response_writer, http.StatusCreated, created_route_deviation_alert_record)
}

func (api_handler *Api_handler) handle_list_route_deviation_alerts(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	if api_handler.route_repository == nil {
		write_json_error(response_writer, http.StatusInternalServerError, "route_repository_not_configured")
		return
	}

	limit_value := 0
	raw_limit_value := strings.TrimSpace(request.URL.Query().Get("limit"))
	if raw_limit_value != "" {
		parsed_limit_value, parse_limit_error := strconv.Atoi(raw_limit_value)
		if parse_limit_error != nil || parsed_limit_value <= 0 {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_limit_query_param")
			return
		}
		limit_value = parsed_limit_value
	}

	route_deviation_alert_list_result, list_alerts_error := api_handler.route_repository.List_route_deviation_alerts(
		request.Context(),
		domain.Route_deviation_alert_list_query{
			Truck_identifier: strings.TrimSpace(request.URL.Query().Get("truck_identifier")),
			Route_identifier: strings.TrimSpace(request.URL.Query().Get("route_identifier")),
			Alert_status:     strings.TrimSpace(request.URL.Query().Get("status")),
			Limit:            limit_value,
		},
	)
	if list_alerts_error != nil {
		http_status_code, error_message := map_route_error_to_http_response(list_alerts_error)
		api_handler.application_logger.Error("failed_to_list_route_deviation_alerts", "error", list_alerts_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, route_deviation_alert_list_result)
}

func map_route_error_to_http_response(route_error error) (int, string) {
	switch route_error.Error() {
	case "invalid_collection_weekday_filter":
		return http.StatusBadRequest, "invalid_collection_weekday_filter"
	case "route_code_is_required":
		return http.StatusBadRequest, "route_code_is_required"
	case "route_name_is_required":
		return http.StatusBadRequest, "route_name_is_required"
	case "zone_name_is_required":
		return http.StatusBadRequest, "zone_name_is_required"
	case "collection_weekday_out_of_range":
		return http.StatusBadRequest, "collection_weekday_out_of_range"
	case "route_code_already_exists":
		return http.StatusConflict, "route_code_already_exists"
	case "no_route_fields_to_update":
		return http.StatusBadRequest, "no_route_fields_to_update"
	case "invalid_route_identifier":
		return http.StatusBadRequest, "invalid_route_identifier"
	case "invalid_bin_identifier":
		return http.StatusBadRequest, "invalid_bin_identifier"
	case "route_stop_order_must_be_positive":
		return http.StatusBadRequest, "route_stop_order_must_be_positive"
	case "duplicate_route_stop_order":
		return http.StatusBadRequest, "duplicate_route_stop_order"
	case "duplicate_bin_identifier_in_route_stops":
		return http.StatusBadRequest, "duplicate_bin_identifier_in_route_stops"
	case "invalid_planned_time":
		return http.StatusBadRequest, "invalid_planned_time"
	case "truck_identifier_is_required":
		return http.StatusBadRequest, "truck_identifier_is_required"
	case "active_route_assignment_not_found_for_truck":
		return http.StatusNotFound, "active_route_assignment_not_found_for_truck"
	case "truck_position_not_found":
		return http.StatusNotFound, "truck_position_not_found"
	case "truck_is_within_route_threshold":
		return http.StatusConflict, "truck_is_within_route_threshold"
	case "route_has_no_stops":
		return http.StatusBadRequest, "route_has_no_stops"
	case "invalid_route_deviation_alert_status_filter":
		return http.StatusBadRequest, "invalid_route_deviation_alert_status_filter"
	case "bin_not_found":
		return http.StatusNotFound, "bin_not_found"
	case "bin_zone_mismatch":
		return http.StatusBadRequest, "bin_zone_mismatch"
	case "route_not_found":
		return http.StatusNotFound, "route_not_found"
	default:
		return http.StatusInternalServerError, "failed_to_process_route_request"
	}
}

func parse_optional_collection_weekday_query_parameter(raw_query_value string) (int, error) {
	if raw_query_value == "" {
		return 0, nil
	}

	parsed_weekday_value, parse_error := strconv.Atoi(raw_query_value)
	if parse_error != nil {
		return 0, parse_error
	}

	if parsed_weekday_value < 1 || parsed_weekday_value > 7 {
		return 0, strconv.ErrSyntax
	}

	return parsed_weekday_value, nil
}

func parse_optional_bool_query_parameter(raw_query_value string) (bool, bool, error) {
	if raw_query_value == "" {
		return false, false, nil
	}

	parsed_bool_value, parse_error := strconv.ParseBool(raw_query_value)
	if parse_error != nil {
		return false, false, parse_error
	}

	return true, parsed_bool_value, nil
}
