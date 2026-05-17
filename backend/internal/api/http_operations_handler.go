package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ecochitas/internal/domain"
)

func (api_handler *Api_handler) handle_ingest_bin_sensor_event(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	var bin_sensor_event_payload struct {
		Bin_identifier    string `json:"bin_identifier"`
		Source_identifier string `json:"source_identifier"`
		Fill_percentage   int    `json:"fill_percentage"`
		Sensor_status     string `json:"sensor_status"`
		Measured_at       string `json:"measured_at"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&bin_sensor_event_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_bin_sensor_event_payload")
		return
	}

	measured_at_value, parse_measured_at_error := parse_optional_utc_timestamp(
		strings.TrimSpace(bin_sensor_event_payload.Measured_at),
	)
	if parse_measured_at_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_measured_at_format")
		return
	}

	bin_sensor_event_ingestion_command := domain.Bin_sensor_event_ingestion_command{
		Bin_identifier:    strings.TrimSpace(bin_sensor_event_payload.Bin_identifier),
		Source_identifier: strings.TrimSpace(bin_sensor_event_payload.Source_identifier),
		Fill_percentage:   bin_sensor_event_payload.Fill_percentage,
		Sensor_status:     strings.TrimSpace(bin_sensor_event_payload.Sensor_status),
		Measured_at:       measured_at_value,
	}

	recorded_sensor_event, ingest_sensor_event_error := api_handler.operations_repository.Ingest_bin_sensor_event(
		query_context,
		bin_sensor_event_ingestion_command,
	)
	if ingest_sensor_event_error != nil {
		http_status_code, error_message := map_operations_error_to_http_response(ingest_sensor_event_error)
		api_handler.application_logger.Error("failed_to_ingest_bin_sensor_event", "error", ingest_sensor_event_error)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	if api_handler.operations_event_publisher != nil {
		publish_event_error := api_handler.operations_event_publisher.Publish_bin_sensor_event(recorded_sensor_event)
		if publish_event_error != nil {
			api_handler.application_logger.Warn("failed_to_publish_bin_sensor_event", "error", publish_event_error)
		}
	}

	write_json_response(response_writer, http.StatusCreated, recorded_sensor_event)
}

func (api_handler *Api_handler) handle_record_driver_collection_event(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	var driver_collection_event_payload struct {
		Route_stop_identifier string `json:"route_stop_identifier"`
		Bin_identifier        string `json:"bin_identifier"`
		Action_type           string `json:"action_type"`
		Evidence_photo_url    string `json:"evidence_photo_url"`
		Action_notes          string `json:"action_notes"`
		Action_at             string `json:"action_at"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&driver_collection_event_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_driver_collection_event_payload")
		return
	}

	action_at_value, parse_action_at_error := parse_optional_utc_timestamp(
		strings.TrimSpace(driver_collection_event_payload.Action_at),
	)
	if parse_action_at_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_action_at_format")
		return
	}

	driver_collection_event_create_command := domain.Driver_collection_event_create_command{
		Route_stop_identifier:         strings.TrimSpace(driver_collection_event_payload.Route_stop_identifier),
		Bin_identifier:                strings.TrimSpace(driver_collection_event_payload.Bin_identifier),
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
		Authenticated_role_name:       strings.TrimSpace(auth_claims.Role_name),
		Authenticated_full_name:       strings.TrimSpace(auth_claims.Full_name),
		Action_type:                   strings.TrimSpace(driver_collection_event_payload.Action_type),
		Evidence_photo_url:            strings.TrimSpace(driver_collection_event_payload.Evidence_photo_url),
		Action_notes:                  strings.TrimSpace(driver_collection_event_payload.Action_notes),
		Action_at:                     action_at_value,
	}

	recorded_driver_collection_event, record_driver_collection_event_error := api_handler.operations_repository.Record_driver_collection_event(
		query_context,
		driver_collection_event_create_command,
	)
	if record_driver_collection_event_error != nil {
		http_status_code, error_message := map_operations_error_to_http_response(record_driver_collection_event_error)
		api_handler.application_logger.Error(
			"failed_to_record_driver_collection_event",
			"error",
			record_driver_collection_event_error,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	if api_handler.operations_event_publisher != nil {
		publish_event_error := api_handler.operations_event_publisher.Publish_driver_collection_event(
			recorded_driver_collection_event,
		)
		if publish_event_error != nil {
			api_handler.application_logger.Warn("failed_to_publish_driver_collection_event", "error", publish_event_error)
		}
	}

	write_json_response(response_writer, http.StatusCreated, recorded_driver_collection_event)
}

func (api_handler *Api_handler) handle_create_route_blockage_report(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	var route_blockage_report_payload struct {
		Route_identifier      string `json:"route_identifier"`
		Route_stop_identifier string `json:"route_stop_identifier"`
		Bin_identifier        string `json:"bin_identifier"`
		Blockage_reason       string `json:"blockage_reason"`
		Evidence_photo_url    string `json:"evidence_photo_url"`
		Severity_level        string `json:"severity_level"`
		Reported_at           string `json:"reported_at"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&route_blockage_report_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_route_blockage_report_payload")
		return
	}

	reported_at_value, parse_reported_at_error := parse_optional_utc_timestamp(
		strings.TrimSpace(route_blockage_report_payload.Reported_at),
	)
	if parse_reported_at_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_reported_at_format")
		return
	}

	route_blockage_report_create_command := domain.Route_blockage_report_create_command{
		Route_identifier:              strings.TrimSpace(route_blockage_report_payload.Route_identifier),
		Route_stop_identifier:         strings.TrimSpace(route_blockage_report_payload.Route_stop_identifier),
		Bin_identifier:                strings.TrimSpace(route_blockage_report_payload.Bin_identifier),
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
		Blockage_reason:               strings.TrimSpace(route_blockage_report_payload.Blockage_reason),
		Evidence_photo_url:            strings.TrimSpace(route_blockage_report_payload.Evidence_photo_url),
		Severity_level:                strings.TrimSpace(route_blockage_report_payload.Severity_level),
		Reported_at:                   reported_at_value,
	}

	created_route_blockage_report, create_route_blockage_report_error := api_handler.operations_repository.Create_route_blockage_report(
		query_context,
		route_blockage_report_create_command,
	)
	if create_route_blockage_report_error != nil {
		http_status_code, error_message := map_operations_error_to_http_response(create_route_blockage_report_error)
		api_handler.application_logger.Error(
			"failed_to_create_route_blockage_report",
			"error",
			create_route_blockage_report_error,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	if api_handler.operations_event_publisher != nil {
		publish_event_error := api_handler.operations_event_publisher.Publish_route_blockage_report_event(
			created_route_blockage_report,
		)
		if publish_event_error != nil {
			api_handler.application_logger.Warn("failed_to_publish_route_blockage_report_event", "error", publish_event_error)
		}
	}

	write_json_response(response_writer, http.StatusCreated, created_route_blockage_report)
}

func (api_handler *Api_handler) handle_list_route_blockage_reports(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	limit_query_value := strings.TrimSpace(request.URL.Query().Get("limit"))
	limit_value := 50
	if limit_query_value != "" {
		parsed_limit_value, parse_limit_error := parse_int_query_parameter(limit_query_value)
		if parse_limit_error != nil {
			write_json_error(response_writer, http.StatusBadRequest, "invalid_limit_query_param")
			return
		}
		limit_value = parsed_limit_value
	}

	route_blockage_report_list_query := domain.Route_blockage_report_list_query{
		Route_identifier:      strings.TrimSpace(request.URL.Query().Get("route_identifier")),
		Route_stop_identifier: strings.TrimSpace(request.URL.Query().Get("route_stop_identifier")),
		Bin_identifier:        strings.TrimSpace(request.URL.Query().Get("bin_identifier")),
		Status_filter:         strings.TrimSpace(request.URL.Query().Get("status")),
		Limit:                 limit_value,
	}

	route_blockage_report_list_result, list_route_blockage_reports_error := api_handler.operations_repository.List_route_blockage_reports(
		query_context,
		route_blockage_report_list_query,
	)
	if list_route_blockage_reports_error != nil {
		http_status_code, error_message := map_operations_error_to_http_response(list_route_blockage_reports_error)
		api_handler.application_logger.Error(
			"failed_to_list_route_blockage_reports",
			"error",
			list_route_blockage_reports_error,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	write_json_response(response_writer, http.StatusOK, route_blockage_report_list_result)
}

func (api_handler *Api_handler) handle_update_route_blockage_report_status(
	response_writer http.ResponseWriter,
	request *http.Request,
) {
	query_context := request.Context()

	auth_claims, has_auth_claims := get_auth_claims_from_request(request)
	if !has_auth_claims {
		write_json_error(response_writer, http.StatusUnauthorized, "auth_claims_not_found")
		return
	}

	blockage_identifier := strings.TrimSpace(request.PathValue("blockage_identifier"))
	if blockage_identifier == "" {
		write_json_error(response_writer, http.StatusBadRequest, "blockage_identifier_path_param_is_required")
		return
	}

	var update_route_blockage_status_payload struct {
		Status           string `json:"status"`
		Resolution_notes string `json:"resolution_notes"`
		Resolved_at      string `json:"resolved_at"`
	}
	decode_payload_error := json.NewDecoder(request.Body).Decode(&update_route_blockage_status_payload)
	if decode_payload_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_update_route_blockage_status_payload")
		return
	}

	resolved_at_value, parse_resolved_at_error := parse_optional_utc_timestamp(
		strings.TrimSpace(update_route_blockage_status_payload.Resolved_at),
	)
	if parse_resolved_at_error != nil {
		write_json_error(response_writer, http.StatusBadRequest, "invalid_resolved_at_format")
		return
	}

	route_blockage_report_status_update_command := domain.Route_blockage_report_status_update_command{
		Blockage_identifier:           blockage_identifier,
		Authenticated_user_identifier: strings.TrimSpace(auth_claims.User_identifier),
		Status:                        strings.TrimSpace(update_route_blockage_status_payload.Status),
		Resolution_notes:              strings.TrimSpace(update_route_blockage_status_payload.Resolution_notes),
		Resolved_at:                   resolved_at_value,
	}

	updated_route_blockage_report, update_route_blockage_report_status_error := api_handler.operations_repository.Update_route_blockage_report_status(
		query_context,
		route_blockage_report_status_update_command,
	)
	if update_route_blockage_report_status_error != nil {
		http_status_code, error_message := map_operations_error_to_http_response(update_route_blockage_report_status_error)
		api_handler.application_logger.Error(
			"failed_to_update_route_blockage_report_status",
			"error",
			update_route_blockage_report_status_error,
		)
		write_json_error(response_writer, http_status_code, error_message)
		return
	}

	if api_handler.operations_event_publisher != nil {
		publish_event_error := api_handler.operations_event_publisher.Publish_route_blockage_report_status_updated_event(
			updated_route_blockage_report,
		)
		if publish_event_error != nil {
			api_handler.application_logger.Warn(
				"failed_to_publish_route_blockage_report_status_updated_event",
				"error",
				publish_event_error,
			)
		}
	}

	write_json_response(response_writer, http.StatusOK, updated_route_blockage_report)
}

func parse_optional_utc_timestamp(raw_timestamp string) (time.Time, error) {
	if raw_timestamp == "" {
		return time.Now().UTC(), nil
	}

	parsed_timestamp, parse_timestamp_error := time.Parse(time.RFC3339, raw_timestamp)
	if parse_timestamp_error != nil {
		return time.Time{}, parse_timestamp_error
	}

	return parsed_timestamp.UTC(), nil
}

func map_operations_error_to_http_response(operation_error error) (int, string) {
	switch operation_error.Error() {
	case "bin_identifier_is_required":
		return http.StatusBadRequest, "bin_identifier_is_required"
	case "fill_percentage_out_of_range":
		return http.StatusBadRequest, "fill_percentage_out_of_range"
	case "invalid_sensor_status":
		return http.StatusBadRequest, "invalid_sensor_status"
	case "invalid_action_type":
		return http.StatusBadRequest, "invalid_action_type"
	case "authenticated_user_identifier_is_required":
		return http.StatusBadRequest, "authenticated_user_identifier_is_required"
	case "invalid_operation_user_role_name":
		return http.StatusBadRequest, "invalid_operation_user_role_name"
	case "blockage_reason_is_required":
		return http.StatusBadRequest, "blockage_reason_is_required"
	case "invalid_severity_level":
		return http.StatusBadRequest, "invalid_severity_level"
	case "invalid_status_filter":
		return http.StatusBadRequest, "invalid_status_filter"
	case "invalid_blockage_status":
		return http.StatusBadRequest, "invalid_blockage_status"
	case "invalid_blockage_identifier":
		return http.StatusBadRequest, "invalid_blockage_identifier"
	case "bin_not_found":
		return http.StatusNotFound, "bin_not_found"
	case "blockage_report_not_found":
		return http.StatusNotFound, "blockage_report_not_found"
	default:
		return http.StatusInternalServerError, "failed_to_process_operation_event"
	}
}

func parse_int_query_parameter(raw_query_value string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(raw_query_value))
}
