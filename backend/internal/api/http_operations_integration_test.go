package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ecochitas/internal/auth"
	"ecochitas/internal/domain"
)

type fake_operations_service struct {
	ingest_bin_sensor_event_response       *domain.Bin_sensor_event_record
	ingest_bin_sensor_event_error          error
	record_driver_event_response           *domain.Driver_collection_event_record
	record_driver_event_error              error
	create_route_blockage_response         *domain.Route_blockage_report_record
	create_route_blockage_error            error
	list_route_blockage_reports_response   *domain.Route_blockage_report_list_result
	list_route_blockage_reports_error      error
	update_route_blockage_response         *domain.Route_blockage_report_record
	update_route_blockage_error            error
	last_ingest_bin_sensor_event_command   *domain.Bin_sensor_event_ingestion_command
	last_record_driver_event_command       *domain.Driver_collection_event_create_command
	last_create_route_blockage_command     *domain.Route_blockage_report_create_command
	last_list_route_blockage_reports_query *domain.Route_blockage_report_list_query
	last_update_route_blockage_command     *domain.Route_blockage_report_status_update_command
}

func (fake_service *fake_operations_service) Ingest_bin_sensor_event(
	application_context context.Context,
	bin_sensor_event_ingestion_command domain.Bin_sensor_event_ingestion_command,
) (*domain.Bin_sensor_event_record, error) {
	_ = application_context
	command_copy := bin_sensor_event_ingestion_command
	fake_service.last_ingest_bin_sensor_event_command = &command_copy
	return fake_service.ingest_bin_sensor_event_response, fake_service.ingest_bin_sensor_event_error
}

func (fake_service *fake_operations_service) Record_driver_collection_event(
	application_context context.Context,
	driver_collection_event_create_command domain.Driver_collection_event_create_command,
) (*domain.Driver_collection_event_record, error) {
	_ = application_context
	command_copy := driver_collection_event_create_command
	fake_service.last_record_driver_event_command = &command_copy
	return fake_service.record_driver_event_response, fake_service.record_driver_event_error
}

func (fake_service *fake_operations_service) Create_route_blockage_report(
	application_context context.Context,
	route_blockage_report_create_command domain.Route_blockage_report_create_command,
) (*domain.Route_blockage_report_record, error) {
	_ = application_context
	command_copy := route_blockage_report_create_command
	fake_service.last_create_route_blockage_command = &command_copy
	return fake_service.create_route_blockage_response, fake_service.create_route_blockage_error
}

func (fake_service *fake_operations_service) List_route_blockage_reports(
	application_context context.Context,
	route_blockage_report_list_query domain.Route_blockage_report_list_query,
) (*domain.Route_blockage_report_list_result, error) {
	_ = application_context
	query_copy := route_blockage_report_list_query
	fake_service.last_list_route_blockage_reports_query = &query_copy
	return fake_service.list_route_blockage_reports_response, fake_service.list_route_blockage_reports_error
}

func (fake_service *fake_operations_service) Update_route_blockage_report_status(
	application_context context.Context,
	route_blockage_report_status_update_command domain.Route_blockage_report_status_update_command,
) (*domain.Route_blockage_report_record, error) {
	_ = application_context
	command_copy := route_blockage_report_status_update_command
	fake_service.last_update_route_blockage_command = &command_copy
	return fake_service.update_route_blockage_response, fake_service.update_route_blockage_error
}

type fake_operations_event_publisher struct {
	published_bin_sensor_event_total                   int
	published_driver_collection_event_total            int
	published_route_blockage_event_total               int
	published_route_blockage_status_update_event_total int
	published_route_deviation_alert_event_total        int
}

func (fake_publisher *fake_operations_event_publisher) Publish_bin_sensor_event(
	bin_sensor_event_record *domain.Bin_sensor_event_record,
) error {
	if bin_sensor_event_record != nil {
		fake_publisher.published_bin_sensor_event_total++
	}
	return nil
}

func (fake_publisher *fake_operations_event_publisher) Publish_driver_collection_event(
	driver_collection_event_record *domain.Driver_collection_event_record,
) error {
	if driver_collection_event_record != nil {
		fake_publisher.published_driver_collection_event_total++
	}
	return nil
}

func (fake_publisher *fake_operations_event_publisher) Publish_route_blockage_report_event(
	route_blockage_report_record *domain.Route_blockage_report_record,
) error {
	if route_blockage_report_record != nil {
		fake_publisher.published_route_blockage_event_total++
	}
	return nil
}

func (fake_publisher *fake_operations_event_publisher) Publish_route_blockage_report_status_updated_event(
	route_blockage_report_record *domain.Route_blockage_report_record,
) error {
	if route_blockage_report_record != nil {
		fake_publisher.published_route_blockage_status_update_event_total++
	}
	return nil
}

func (fake_publisher *fake_operations_event_publisher) Publish_route_deviation_alert_event(
	route_deviation_alert_record *domain.Route_deviation_alert_record,
) error {
	if route_deviation_alert_record != nil {
		fake_publisher.published_route_deviation_alert_event_total++
	}
	return nil
}

func Test_operations_endpoints_require_authentication(t *testing.T) {
	api_http_multiplexer, _, _ := create_operations_test_server(t)

	request_body := bytes.NewBufferString(`{"bin_identifier":"1f1e9eb7-5f6e-4c65-b84e-b9650ce4d3ee","fill_percentage":10,"sensor_status":"online"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/bins/sensor-events", request_body)
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "authorization_header_is_required" {
		t.Fatalf("expected authorization_header_is_required, got %s", error_response.Message)
	}
}

func Test_operations_endpoints_validate_role_permissions(t *testing.T) {
	api_http_multiplexer, api_handler, _ := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "citizen")

	request_body := bytes.NewBufferString(`{"bin_identifier":"1f1e9eb7-5f6e-4c65-b84e-b9650ce4d3ee","fill_percentage":10,"sensor_status":"online"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/bins/sensor-events", request_body)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "insufficient_permissions" {
		t.Fatalf("expected insufficient_permissions, got %s", error_response.Message)
	}
}

func Test_ingest_bin_sensor_event_validates_timestamp_before_service_call(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "driver")

	request_body := bytes.NewBufferString(`{"bin_identifier":"1f1e9eb7-5f6e-4c65-b84e-b9650ce4d3ee","fill_percentage":10,"sensor_status":"online","measured_at":"invalid-date"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/bins/sensor-events", request_body)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "invalid_measured_at_format" {
		t.Fatalf("expected invalid_measured_at_format, got %s", error_response.Message)
	}

	if fake_service.last_ingest_bin_sensor_event_command != nil {
		t.Fatalf("service should not be called when measured_at is invalid")
	}
}

func Test_list_route_blockage_reports_happy_path(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "admin")
	fake_service.list_route_blockage_reports_response = &domain.Route_blockage_report_list_result{
		Items: []domain.Route_blockage_report_record{
			{
				Blockage_identifier:         "6b1bbbd5-bf32-4727-83f0-2f8657da98f2",
				Reported_by_user_identifier: "driver-demo-001",
				Blockage_reason:             "vehicle blocking route",
				Severity_level:              "medium",
				Status:                      "open",
				Reported_at:                 time.Now().UTC(),
				Created_at:                  time.Now().UTC(),
				Updated_at:                  time.Now().UTC(),
			},
		},
		Total: 1,
	}

	http_request := httptest.NewRequest(
		http.MethodGet,
		"/v1/driver/route-blockages?status=open&limit=10",
		nil,
	)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, http_response_recorder.Code)
	}

	if fake_service.last_list_route_blockage_reports_query == nil {
		t.Fatalf("expected service query to be captured")
	}
	if fake_service.last_list_route_blockage_reports_query.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", fake_service.last_list_route_blockage_reports_query.Limit)
	}
	if fake_service.last_list_route_blockage_reports_query.Status_filter != "open" {
		t.Fatalf("expected status filter open, got %s", fake_service.last_list_route_blockage_reports_query.Status_filter)
	}
}

func Test_ingest_bin_sensor_event_happy_path_publishes_event(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "driver")
	fake_service.ingest_bin_sensor_event_response = &domain.Bin_sensor_event_record{
		Sensor_event_identifier: "1",
		Bin_identifier:          "1f1e9eb7-5f6e-4c65-b84e-b9650ce4d3ee",
		Source_identifier:       "sensor-demo-001",
		Fill_percentage:         90,
		Sensor_status:           "online",
		Bin_status:              "full",
		Measured_at:             time.Now().UTC(),
		Recorded_at:             time.Now().UTC(),
	}

	request_body := bytes.NewBufferString(`{"bin_identifier":"1f1e9eb7-5f6e-4c65-b84e-b9650ce4d3ee","fill_percentage":90,"sensor_status":"online","source_identifier":"sensor-demo-001","measured_at":"2026-05-17T10:30:00Z"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/bins/sensor-events", request_body)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, http_response_recorder.Code)
	}

	typed_publisher, is_expected_publisher_type := api_handler.operations_event_publisher.(*fake_operations_event_publisher)
	if !is_expected_publisher_type {
		t.Fatalf("expected fake publisher type")
	}
	if typed_publisher.published_bin_sensor_event_total != 1 {
		t.Fatalf(
			"expected 1 published bin sensor event, got %d",
			typed_publisher.published_bin_sensor_event_total,
		)
	}
}

func create_operations_test_server(
	t *testing.T,
) (*http.ServeMux, *Api_handler, *fake_operations_service) {
	t.Helper()

	jwt_authenticator, create_authenticator_error := auth.New_jwt_authenticator(
		"test_signing_key",
		"ecochitas_backend",
		"ecochitas_api",
		60,
	)
	if create_authenticator_error != nil {
		t.Fatalf("failed to create jwt authenticator: %v", create_authenticator_error)
	}

	fake_service := &fake_operations_service{}
	fake_publisher := &fake_operations_event_publisher{}
	discard_logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	api_handler := New_api_handler(
		nil,
		nil,
		nil,
		nil,
		fake_service,
		nil,
		fake_publisher,
		jwt_authenticator,
		true,
		discard_logger,
	)
	api_http_multiplexer := http.NewServeMux()
	api_handler.Register_routes(api_http_multiplexer)

	return api_http_multiplexer, api_handler, fake_service
}

func must_issue_access_token(t *testing.T, api_handler *Api_handler, role_name string) string {
	t.Helper()

	access_token, _, generate_token_error := api_handler.jwt_authenticator.Generate_access_token(
		auth.Generate_access_token_command{
			User_identifier: "driver-demo-001",
			Role_name:       role_name,
			Full_name:       "Driver Demo",
		},
	)
	if generate_token_error != nil {
		t.Fatalf("failed to generate token: %v", generate_token_error)
	}

	return access_token
}

func decode_error_response(t *testing.T, response_body *bytes.Buffer) Error_response {
	t.Helper()

	var error_response Error_response
	decode_error := json.NewDecoder(response_body).Decode(&error_response)
	if decode_error != nil {
		t.Fatalf("failed to decode error response: %v", decode_error)
	}
	if error_response.Message == "" {
		t.Fatalf("error response message should not be empty")
	}
	return error_response
}

func Test_list_route_blockage_reports_rejects_invalid_limit_query_param(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "admin")

	http_request := httptest.NewRequest(
		http.MethodGet,
		"/v1/driver/route-blockages?limit=invalid",
		nil,
	)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "invalid_limit_query_param" {
		t.Fatalf("expected invalid_limit_query_param, got %s", error_response.Message)
	}

	if fake_service.last_list_route_blockage_reports_query != nil {
		t.Fatalf("service should not be called with invalid limit query param")
	}
}

func Test_map_operations_error_to_http_response(t *testing.T) {
	http_status_code, error_message := map_operations_error_to_http_response(fmt.Errorf("invalid_status_filter"))
	if http_status_code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, http_status_code)
	}
	if error_message != "invalid_status_filter" {
		t.Fatalf("expected invalid_status_filter, got %s", error_message)
	}
}

func Test_update_route_blockage_report_status_rejects_invalid_resolved_at(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "driver")

	request_body := bytes.NewBufferString(`{"status":"resolved","resolution_notes":"fixed","resolved_at":"invalid-date"}`)
	http_request := httptest.NewRequest(
		http.MethodPatch,
		"/v1/driver/route-blockages/6b1bbbd5-bf32-4727-83f0-2f8657da98f2",
		request_body,
	)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "invalid_resolved_at_format" {
		t.Fatalf("expected invalid_resolved_at_format, got %s", error_response.Message)
	}
	if fake_service.last_update_route_blockage_command != nil {
		t.Fatalf("service should not be called when resolved_at is invalid")
	}
}

func Test_update_route_blockage_report_status_happy_path_publishes_event(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_operations_test_server(t)
	access_token := must_issue_access_token(t, api_handler, "admin")
	fake_service.update_route_blockage_response = &domain.Route_blockage_report_record{
		Blockage_identifier:         "6b1bbbd5-bf32-4727-83f0-2f8657da98f2",
		Reported_by_user_identifier: "driver-demo-001",
		Blockage_reason:             "vehicle blocking route",
		Severity_level:              "medium",
		Status:                      "resolved",
		Reported_at:                 time.Now().UTC(),
		Resolved_at:                 pointer_to_time(time.Now().UTC()),
		Resolution_notes:            "access restored",
		Created_at:                  time.Now().UTC(),
		Updated_at:                  time.Now().UTC(),
	}

	request_body := bytes.NewBufferString(`{"status":"resolved","resolution_notes":"access restored","resolved_at":"2026-05-17T12:30:00Z"}`)
	http_request := httptest.NewRequest(
		http.MethodPatch,
		"/v1/driver/route-blockages/6b1bbbd5-bf32-4727-83f0-2f8657da98f2",
		request_body,
	)
	http_request.Header.Set("Authorization", "Bearer "+access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, http_response_recorder.Code)
	}

	if fake_service.last_update_route_blockage_command == nil {
		t.Fatalf("expected update command to be sent to service")
	}
	if fake_service.last_update_route_blockage_command.Blockage_identifier != "6b1bbbd5-bf32-4727-83f0-2f8657da98f2" {
		t.Fatalf("unexpected blockage identifier in update command")
	}
	if fake_service.last_update_route_blockage_command.Status != "resolved" {
		t.Fatalf("expected status resolved, got %s", fake_service.last_update_route_blockage_command.Status)
	}

	typed_publisher, is_expected_publisher_type := api_handler.operations_event_publisher.(*fake_operations_event_publisher)
	if !is_expected_publisher_type {
		t.Fatalf("expected fake publisher type")
	}
	if typed_publisher.published_route_blockage_status_update_event_total != 1 {
		t.Fatalf(
			"expected 1 published route blockage status update event, got %d",
			typed_publisher.published_route_blockage_status_update_event_total,
		)
	}
}

func pointer_to_time(time_value time.Time) *time.Time {
	return &time_value
}
