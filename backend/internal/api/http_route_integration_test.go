package api

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ecochitas/internal/auth"
	"ecochitas/internal/domain"
)

type fake_route_service struct {
	list_collection_routes_response          []domain.Collection_route_view
	list_collection_routes_error             error
	create_collection_route_response         *domain.Collection_route_view
	create_collection_route_error            error
	update_collection_route_response         *domain.Collection_route_view
	update_collection_route_error            error
	list_route_stops_response                []domain.Route_stop_view
	list_route_stops_error                   error
	sync_route_stops_response                []domain.Route_stop_view
	sync_route_stops_error                   error
	list_collection_route_revisions_response *domain.Collection_route_revision_list_result
	list_collection_route_revisions_error    error
	get_truck_route_deviation_response       *domain.Truck_route_deviation_view
	get_truck_route_deviation_error          error
	create_truck_route_assignment_response   *domain.Truck_route_assignment_view
	create_truck_route_assignment_error      error
	list_truck_route_assignments_response    *domain.Truck_route_assignment_list_result
	list_truck_route_assignments_error       error
	create_route_deviation_alert_response    *domain.Route_deviation_alert_record
	create_route_deviation_alert_error       error
	list_route_deviation_alerts_response     *domain.Route_deviation_alert_list_result
	list_route_deviation_alerts_error        error
	last_create_collection_route_command     *domain.Collection_route_create_command
	last_update_collection_route_command     *domain.Collection_route_update_command
	last_sync_route_stops_command            *domain.Route_stop_sync_command
	last_list_revisions_query                *domain.Collection_route_revision_list_query
	last_truck_route_deviation_query         *domain.Truck_route_deviation_query
	last_create_truck_route_assignment       *domain.Truck_route_assignment_create_command
	last_list_truck_route_assignments_query  *domain.Truck_route_assignment_list_query
	last_create_route_deviation_alert        *domain.Route_deviation_alert_create_command
	last_list_route_deviation_alerts_query   *domain.Route_deviation_alert_list_query
}

func (fake_service *fake_route_service) List_collection_routes(
	application_context context.Context,
	collection_route_list_query domain.Collection_route_list_query,
) ([]domain.Collection_route_view, error) {
	_ = application_context
	_ = collection_route_list_query
	return fake_service.list_collection_routes_response, fake_service.list_collection_routes_error
}

func (fake_service *fake_route_service) Create_collection_route(
	application_context context.Context,
	collection_route_create_command domain.Collection_route_create_command,
) (*domain.Collection_route_view, error) {
	_ = application_context
	command_copy := collection_route_create_command
	fake_service.last_create_collection_route_command = &command_copy
	return fake_service.create_collection_route_response, fake_service.create_collection_route_error
}

func (fake_service *fake_route_service) Update_collection_route(
	application_context context.Context,
	collection_route_update_command domain.Collection_route_update_command,
) (*domain.Collection_route_view, error) {
	_ = application_context
	command_copy := collection_route_update_command
	fake_service.last_update_collection_route_command = &command_copy
	return fake_service.update_collection_route_response, fake_service.update_collection_route_error
}

func (fake_service *fake_route_service) List_route_stops_by_route_identifier(
	application_context context.Context,
	route_identifier string,
) ([]domain.Route_stop_view, error) {
	_ = application_context
	_ = route_identifier
	return fake_service.list_route_stops_response, fake_service.list_route_stops_error
}

func (fake_service *fake_route_service) Sync_route_stops(
	application_context context.Context,
	route_stop_sync_command domain.Route_stop_sync_command,
) ([]domain.Route_stop_view, error) {
	_ = application_context
	command_copy := route_stop_sync_command
	fake_service.last_sync_route_stops_command = &command_copy
	return fake_service.sync_route_stops_response, fake_service.sync_route_stops_error
}

func (fake_service *fake_route_service) List_collection_route_revisions(
	application_context context.Context,
	collection_route_revision_list_query domain.Collection_route_revision_list_query,
) (*domain.Collection_route_revision_list_result, error) {
	_ = application_context
	query_copy := collection_route_revision_list_query
	fake_service.last_list_revisions_query = &query_copy
	return fake_service.list_collection_route_revisions_response, fake_service.list_collection_route_revisions_error
}

func (fake_service *fake_route_service) Get_truck_route_deviation(
	application_context context.Context,
	truck_route_deviation_query domain.Truck_route_deviation_query,
) (*domain.Truck_route_deviation_view, error) {
	_ = application_context
	query_copy := truck_route_deviation_query
	fake_service.last_truck_route_deviation_query = &query_copy
	return fake_service.get_truck_route_deviation_response, fake_service.get_truck_route_deviation_error
}

func (fake_service *fake_route_service) Create_truck_route_assignment(
	application_context context.Context,
	truck_route_assignment_create_command domain.Truck_route_assignment_create_command,
) (*domain.Truck_route_assignment_view, error) {
	_ = application_context
	command_copy := truck_route_assignment_create_command
	fake_service.last_create_truck_route_assignment = &command_copy
	return fake_service.create_truck_route_assignment_response, fake_service.create_truck_route_assignment_error
}

func (fake_service *fake_route_service) List_truck_route_assignments(
	application_context context.Context,
	truck_route_assignment_list_query domain.Truck_route_assignment_list_query,
) (*domain.Truck_route_assignment_list_result, error) {
	_ = application_context
	query_copy := truck_route_assignment_list_query
	fake_service.last_list_truck_route_assignments_query = &query_copy
	return fake_service.list_truck_route_assignments_response, fake_service.list_truck_route_assignments_error
}

func (fake_service *fake_route_service) Create_route_deviation_alert(
	application_context context.Context,
	route_deviation_alert_create_command domain.Route_deviation_alert_create_command,
) (*domain.Route_deviation_alert_record, error) {
	_ = application_context
	command_copy := route_deviation_alert_create_command
	fake_service.last_create_route_deviation_alert = &command_copy
	return fake_service.create_route_deviation_alert_response, fake_service.create_route_deviation_alert_error
}

func (fake_service *fake_route_service) List_route_deviation_alerts(
	application_context context.Context,
	route_deviation_alert_list_query domain.Route_deviation_alert_list_query,
) (*domain.Route_deviation_alert_list_result, error) {
	_ = application_context
	query_copy := route_deviation_alert_list_query
	fake_service.last_list_route_deviation_alerts_query = &query_copy
	return fake_service.list_route_deviation_alerts_response, fake_service.list_route_deviation_alerts_error
}

func Test_admin_collection_route_create_requires_admin_role(t *testing.T) {
	api_http_multiplexer, api_handler, _ := create_route_test_server(t)
	driver_access_token := must_issue_access_token(t, api_handler, "driver")

	request_body := bytes.NewBufferString(`{"route_code":"RUTA-001","route_name":"Ruta Centro","zone_name":"Cochabamba Centro","collection_weekday":2}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/admin/collection-routes", request_body)
	http_request.Header.Set("Authorization", "Bearer "+driver_access_token)
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

func Test_admin_collection_route_create_happy_path(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_route_test_server(t)
	admin_access_token := must_issue_access_token(t, api_handler, "admin")

	fake_service.create_collection_route_response = &domain.Collection_route_view{
		Route_identifier:   "f33403f5-89c8-4f2f-a26f-c4958ac6e577",
		Route_code:         "RUTA-001",
		Route_name:         "Ruta Centro",
		Zone_name:          "Cochabamba Centro",
		Collection_weekday: 2,
		Is_active:          true,
		Stop_total:         0,
		Path_coordinates:   []domain.Route_path_coordinate{},
	}

	request_body := bytes.NewBufferString(`{"route_code":"RUTA-001","route_name":"Ruta Centro","zone_name":"Cochabamba Centro","collection_weekday":2}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/admin/collection-routes", request_body)
	http_request.Header.Set("Authorization", "Bearer "+admin_access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, http_response_recorder.Code)
	}

	if fake_service.last_create_collection_route_command == nil {
		t.Fatalf("expected create command to be sent to service")
	}
	if fake_service.last_create_collection_route_command.Authenticated_user_identifier != "driver-demo-001" {
		t.Fatalf("unexpected authenticated user identifier")
	}
	if fake_service.last_create_collection_route_command.Collection_weekday != 2 {
		t.Fatalf("expected collection_weekday 2, got %d", fake_service.last_create_collection_route_command.Collection_weekday)
	}
}

func Test_list_collection_route_revisions_rejects_invalid_limit(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_route_test_server(t)
	admin_access_token := must_issue_access_token(t, api_handler, "admin")

	http_request := httptest.NewRequest(
		http.MethodGet,
		"/v1/admin/collection-routes/f33403f5-89c8-4f2f-a26f-c4958ac6e577/revisions?limit=invalid",
		nil,
	)
	http_request.Header.Set("Authorization", "Bearer "+admin_access_token)
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, http_response_recorder.Code)
	}

	error_response := decode_error_response(t, http_response_recorder.Body)
	if error_response.Message != "invalid_limit_query_param" {
		t.Fatalf("expected invalid_limit_query_param, got %s", error_response.Message)
	}

	if fake_service.last_list_revisions_query != nil {
		t.Fatalf("service should not be called with invalid limit")
	}
}

func Test_get_truck_route_deviation_happy_path(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_route_test_server(t)
	driver_access_token := must_issue_access_token(t, api_handler, "driver")

	fake_service.get_truck_route_deviation_response = &domain.Truck_route_deviation_view{
		Truck_identifier:              "TRUCK-001",
		Route_identifier:              "f33403f5-89c8-4f2f-a26f-c4958ac6e577",
		Captured_at:                   time.Now().UTC(),
		Truck_latitude:                -17.3921,
		Truck_longitude:               -66.1561,
		Nearest_route_stop_identifier: "d9b501fd-6bcf-42af-a0df-d7ce12e94ec5",
		Nearest_stop_order:            1,
		Nearest_bin_identifier:        "3e75dfde-6f17-4b31-a24b-097c541db923",
		Nearest_bin_code:              "BIN-001",
		Nearest_bin_latitude:          -17.3919,
		Nearest_bin_longitude:         -66.1564,
		Distance_to_route_meters:      43.2,
		Deviation_threshold_meters:    100,
		Is_off_route:                  false,
	}

	http_request := httptest.NewRequest(
		http.MethodGet,
		"/v1/driver/route-deviation?truck_identifier=TRUCK-001&route_identifier=f33403f5-89c8-4f2f-a26f-c4958ac6e577&deviation_threshold_meters=100",
		nil,
	)
	http_request.Header.Set("Authorization", "Bearer "+driver_access_token)
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, http_response_recorder.Code)
	}

	if fake_service.last_truck_route_deviation_query == nil {
		t.Fatalf("expected route deviation query to be sent to service")
	}
	if fake_service.last_truck_route_deviation_query.Deviation_threshold_meters != 100 {
		t.Fatalf("expected deviation_threshold_meters 100, got %v", fake_service.last_truck_route_deviation_query.Deviation_threshold_meters)
	}
}

func Test_create_truck_route_assignment_happy_path(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_route_test_server(t)
	admin_access_token := must_issue_access_token(t, api_handler, "admin")
	fake_service.create_truck_route_assignment_response = &domain.Truck_route_assignment_view{
		Assignment_identifier: "57c8b311-5462-47bd-a12d-9be2294a36a0",
		Truck_identifier:      "TRUCK-001",
		Route_identifier:      "f33403f5-89c8-4f2f-a26f-c4958ac6e577",
		Route_code:            "RUTA-001",
		Route_name:            "Ruta Centro",
		Zone_name:             "Cochabamba Centro",
		Collection_weekday:    2,
		Is_active:             true,
		Assigned_at:           time.Now().UTC(),
	}

	request_body := bytes.NewBufferString(`{"truck_identifier":"TRUCK-001","route_identifier":"f33403f5-89c8-4f2f-a26f-c4958ac6e577","assignment_notes":"turno manana"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/admin/truck-route-assignments", request_body)
	http_request.Header.Set("Authorization", "Bearer "+admin_access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, http_response_recorder.Code)
	}
	if fake_service.last_create_truck_route_assignment == nil {
		t.Fatalf("expected create truck route assignment command")
	}
	if fake_service.last_create_truck_route_assignment.Assigned_by_user_identifier != "driver-demo-001" {
		t.Fatalf("unexpected assigned_by_user_identifier")
	}
}

func Test_create_route_deviation_alert_happy_path(t *testing.T) {
	api_http_multiplexer, api_handler, fake_service := create_route_test_server(t)
	driver_access_token := must_issue_access_token(t, api_handler, "driver")
	fake_service.create_route_deviation_alert_response = &domain.Route_deviation_alert_record{
		Alert_identifier:           "e66d912e-ab5d-4543-a2f5-922c3f830a4b",
		Truck_identifier:           "TRUCK-001",
		Route_identifier:           "f33403f5-89c8-4f2f-a26f-c4958ac6e577",
		Distance_to_route_meters:   350,
		Deviation_threshold_meters: 200,
		Severity_level:             "medium",
		Alert_status:               "open",
		Detected_at:                time.Now().UTC(),
		Created_at:                 time.Now().UTC(),
		Metadata_payload:           map[string]any{"source": "test"},
	}

	request_body := bytes.NewBufferString(`{"truck_identifier":"TRUCK-001","deviation_threshold_meters":200,"alert_notes":"desvio detectado en prueba"}`)
	http_request := httptest.NewRequest(http.MethodPost, "/v1/driver/route-deviation-alerts", request_body)
	http_request.Header.Set("Authorization", "Bearer "+driver_access_token)
	http_request.Header.Set("Content-Type", "application/json")
	http_response_recorder := httptest.NewRecorder()

	api_http_multiplexer.ServeHTTP(http_response_recorder, http_request)
	if http_response_recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, http_response_recorder.Code)
	}
	if fake_service.last_create_route_deviation_alert == nil {
		t.Fatalf("expected create route deviation alert command")
	}
	if fake_service.last_create_route_deviation_alert.Triggered_by_user_identifier != "driver-demo-001" {
		t.Fatalf("unexpected triggered_by_user_identifier")
	}
}

func create_route_test_server(
	t *testing.T,
) (*http.ServeMux, *Api_handler, *fake_route_service) {
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

	fake_service := &fake_route_service{}
	discard_logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	api_handler := New_api_handler(
		nil,
		nil,
		fake_service,
		nil,
		nil,
		nil,
		nil,
		jwt_authenticator,
		true,
		discard_logger,
	)
	api_http_multiplexer := http.NewServeMux()
	api_handler.Register_routes(api_http_multiplexer)

	return api_http_multiplexer, api_handler, fake_service
}
