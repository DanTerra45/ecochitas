package api

import (
	"context"

	"ecochitas/internal/domain"
)

type operations_service interface {
	Ingest_bin_sensor_event(
		application_context context.Context,
		bin_sensor_event_ingestion_command domain.Bin_sensor_event_ingestion_command,
	) (*domain.Bin_sensor_event_record, error)
	Record_driver_collection_event(
		application_context context.Context,
		driver_collection_event_create_command domain.Driver_collection_event_create_command,
	) (*domain.Driver_collection_event_record, error)
	Create_route_blockage_report(
		application_context context.Context,
		route_blockage_report_create_command domain.Route_blockage_report_create_command,
	) (*domain.Route_blockage_report_record, error)
	List_route_blockage_reports(
		application_context context.Context,
		route_blockage_report_list_query domain.Route_blockage_report_list_query,
	) (*domain.Route_blockage_report_list_result, error)
	Update_route_blockage_report_status(
		application_context context.Context,
		route_blockage_report_status_update_command domain.Route_blockage_report_status_update_command,
	) (*domain.Route_blockage_report_record, error)
}

type operations_event_publisher_service interface {
	Publish_bin_sensor_event(bin_sensor_event_record *domain.Bin_sensor_event_record) error
	Publish_driver_collection_event(driver_collection_event_record *domain.Driver_collection_event_record) error
	Publish_route_blockage_report_event(route_blockage_report_record *domain.Route_blockage_report_record) error
	Publish_route_blockage_report_status_updated_event(route_blockage_report_record *domain.Route_blockage_report_record) error
	Publish_route_deviation_alert_event(route_deviation_alert_record *domain.Route_deviation_alert_record) error
}
