package realtime

import (
	"encoding/json"
	"fmt"
	"time"

	"ecochitas/internal/domain"

	"github.com/nats-io/nats.go"
)

type Operations_event_publisher struct {
	nats_connection                     *nats.Conn
	bin_sensor_event_subject            string
	driver_collection_event_subject     string
	route_blockage_report_event_subject string
	route_deviation_alert_event_subject string
}

type operations_event_payload_envelope struct {
	Event_name   string    `json:"event_name"`
	Published_at time.Time `json:"published_at"`
	Payload      any       `json:"payload"`
}

func New_operations_event_publisher(
	nats_connection *nats.Conn,
	bin_sensor_event_subject string,
	driver_collection_event_subject string,
	route_blockage_report_event_subject string,
	route_deviation_alert_event_subject string,
) *Operations_event_publisher {
	return &Operations_event_publisher{
		nats_connection:                     nats_connection,
		bin_sensor_event_subject:            bin_sensor_event_subject,
		driver_collection_event_subject:     driver_collection_event_subject,
		route_blockage_report_event_subject: route_blockage_report_event_subject,
		route_deviation_alert_event_subject: route_deviation_alert_event_subject,
	}
}

func (operations_event_publisher *Operations_event_publisher) Publish_bin_sensor_event(
	bin_sensor_event_record *domain.Bin_sensor_event_record,
) error {
	return operations_event_publisher.publish_event(
		operations_event_publisher.bin_sensor_event_subject,
		"bin_sensor_event_recorded",
		bin_sensor_event_record,
	)
}

func (operations_event_publisher *Operations_event_publisher) Publish_driver_collection_event(
	driver_collection_event_record *domain.Driver_collection_event_record,
) error {
	return operations_event_publisher.publish_event(
		operations_event_publisher.driver_collection_event_subject,
		"driver_collection_event_recorded",
		driver_collection_event_record,
	)
}

func (operations_event_publisher *Operations_event_publisher) Publish_route_blockage_report_event(
	route_blockage_report_record *domain.Route_blockage_report_record,
) error {
	return operations_event_publisher.publish_event(
		operations_event_publisher.route_blockage_report_event_subject,
		"route_blockage_report_recorded",
		route_blockage_report_record,
	)
}

func (operations_event_publisher *Operations_event_publisher) Publish_route_blockage_report_status_updated_event(
	route_blockage_report_record *domain.Route_blockage_report_record,
) error {
	return operations_event_publisher.publish_event(
		operations_event_publisher.route_blockage_report_event_subject,
		"route_blockage_report_status_updated",
		route_blockage_report_record,
	)
}

func (operations_event_publisher *Operations_event_publisher) Publish_route_deviation_alert_event(
	route_deviation_alert_record *domain.Route_deviation_alert_record,
) error {
	return operations_event_publisher.publish_event(
		operations_event_publisher.route_deviation_alert_event_subject,
		"route_deviation_alert_recorded",
		route_deviation_alert_record,
	)
}

func (operations_event_publisher *Operations_event_publisher) publish_event(
	event_subject string,
	event_name string,
	event_payload any,
) error {
	if operations_event_publisher == nil || operations_event_publisher.nats_connection == nil {
		return nil
	}

	serialized_payload, serialize_payload_error := json.Marshal(operations_event_payload_envelope{
		Event_name:   event_name,
		Published_at: time.Now().UTC(),
		Payload:      event_payload,
	})
	if serialize_payload_error != nil {
		return fmt.Errorf("failed_to_serialize_operations_event_payload: %w", serialize_payload_error)
	}

	publish_message_error := operations_event_publisher.nats_connection.Publish(
		event_subject,
		serialized_payload,
	)
	if publish_message_error != nil {
		return fmt.Errorf("failed_to_publish_operations_event: %w", publish_message_error)
	}

	return nil
}
