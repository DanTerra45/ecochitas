package api

import (
	"context"

	"ecochitas/internal/domain"
)

type route_service interface {
	List_collection_routes(
		application_context context.Context,
		collection_route_list_query domain.Collection_route_list_query,
	) ([]domain.Collection_route_view, error)
	Create_collection_route(
		application_context context.Context,
		collection_route_create_command domain.Collection_route_create_command,
	) (*domain.Collection_route_view, error)
	Create_demo_route(
		application_context context.Context,
		demo_route_create_command domain.Demo_route_create_command,
	) (*domain.Collection_route_view, error)
	Update_collection_route(
		application_context context.Context,
		collection_route_update_command domain.Collection_route_update_command,
	) (*domain.Collection_route_view, error)
	List_route_stops_by_route_identifier(
		application_context context.Context,
		route_identifier string,
	) ([]domain.Route_stop_view, error)
	Sync_route_stops(
		application_context context.Context,
		route_stop_sync_command domain.Route_stop_sync_command,
	) ([]domain.Route_stop_view, error)
	List_collection_route_revisions(
		application_context context.Context,
		collection_route_revision_list_query domain.Collection_route_revision_list_query,
	) (*domain.Collection_route_revision_list_result, error)
	Get_truck_route_deviation(
		application_context context.Context,
		truck_route_deviation_query domain.Truck_route_deviation_query,
	) (*domain.Truck_route_deviation_view, error)
	Create_truck_route_assignment(
		application_context context.Context,
		truck_route_assignment_create_command domain.Truck_route_assignment_create_command,
	) (*domain.Truck_route_assignment_view, error)
	List_truck_route_assignments(
		application_context context.Context,
		truck_route_assignment_list_query domain.Truck_route_assignment_list_query,
	) (*domain.Truck_route_assignment_list_result, error)
	Create_route_deviation_alert(
		application_context context.Context,
		route_deviation_alert_create_command domain.Route_deviation_alert_create_command,
	) (*domain.Route_deviation_alert_record, error)
	List_route_deviation_alerts(
		application_context context.Context,
		route_deviation_alert_list_query domain.Route_deviation_alert_list_query,
	) (*domain.Route_deviation_alert_list_result, error)
}
