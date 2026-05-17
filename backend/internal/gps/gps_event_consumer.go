package gps

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"ecochitas/internal/domain"
	"ecochitas/internal/realtime"
	"ecochitas/internal/storage"

	"github.com/nats-io/nats.go"
)

type Gps_event_consumer struct {
	nats_connection           *nats.Conn
	gps_subject               string
	queue_group_name          string
	truck_position_repository *storage.Truck_position_repository
	truck_position_stream     *realtime.Truck_position_stream
	application_logger        *slog.Logger
	subscription_handler      *nats.Subscription
}

func New_gps_event_consumer(
	nats_connection *nats.Conn,
	gps_subject string,
	queue_group_name string,
	truck_position_repository *storage.Truck_position_repository,
	truck_position_stream *realtime.Truck_position_stream,
	application_logger *slog.Logger,
) *Gps_event_consumer {
	return &Gps_event_consumer{
		nats_connection:           nats_connection,
		gps_subject:               gps_subject,
		queue_group_name:          queue_group_name,
		truck_position_repository: truck_position_repository,
		truck_position_stream:     truck_position_stream,
		application_logger:        application_logger,
	}
}

func (gps_event_consumer *Gps_event_consumer) Start() error {
	nats_subscription, subscribe_error := gps_event_consumer.nats_connection.QueueSubscribe(
		gps_event_consumer.gps_subject,
		gps_event_consumer.queue_group_name,
		gps_event_consumer.process_gps_message,
	)
	if subscribe_error != nil {
		return fmt.Errorf("failed_to_subscribe_gps_subject: %w", subscribe_error)
	}

	gps_event_consumer.subscription_handler = nats_subscription
	return nil
}

func (gps_event_consumer *Gps_event_consumer) Stop() error {
	if gps_event_consumer.subscription_handler == nil {
		return nil
	}

	drain_subscription_error := gps_event_consumer.subscription_handler.Drain()
	if drain_subscription_error != nil {
		return fmt.Errorf("failed_to_drain_gps_subscription: %w", drain_subscription_error)
	}

	return nil
}

func (gps_event_consumer *Gps_event_consumer) process_gps_message(nats_message *nats.Msg) {
	var gps_location_event domain.Gps_location_event

	unmarshal_payload_error := json.Unmarshal(nats_message.Data, &gps_location_event)
	if unmarshal_payload_error != nil {
		gps_event_consumer.application_logger.Error(
			"invalid_gps_event_payload",
			"error",
			unmarshal_payload_error,
			"raw_payload",
			string(nats_message.Data),
		)
		return
	}

	validate_event_error := validate_gps_event(gps_location_event)
	if validate_event_error != nil {
		gps_event_consumer.application_logger.Error(
			"invalid_gps_event_data",
			"error",
			validate_event_error,
			"truck_identifier",
			gps_location_event.Truck_identifier,
		)
		return
	}

	insert_context, cancel_context := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel_context()

	save_position_error := gps_event_consumer.truck_position_repository.Save_truck_position(
		insert_context,
		gps_location_event,
	)
	if save_position_error != nil {
		gps_event_consumer.application_logger.Error(
			"failed_to_persist_gps_event",
			"error",
			save_position_error,
			"truck_identifier",
			gps_location_event.Truck_identifier,
		)
		return
	}

	truck_latest_position := domain.Truck_latest_position{
		Truck_identifier: gps_location_event.Truck_identifier,
		Latitude:         gps_location_event.Latitude,
		Longitude:        gps_location_event.Longitude,
		Speed_kmh:        gps_location_event.Speed_kmh,
		Heading_degrees:  gps_location_event.Heading_degrees,
		Captured_at:      gps_location_event.Captured_at,
		Received_at:      time.Now().UTC(),
	}
	gps_event_consumer.truck_position_stream.Publish(truck_latest_position)

	gps_event_consumer.application_logger.Info(
		"gps_event_processed",
		"truck_identifier",
		gps_location_event.Truck_identifier,
		"captured_at",
		gps_location_event.Captured_at,
	)
}

func validate_gps_event(gps_location_event domain.Gps_location_event) error {
	trimmed_truck_identifier := strings.TrimSpace(gps_location_event.Truck_identifier)
	if trimmed_truck_identifier == "" {
		return fmt.Errorf("truck_identifier_is_required")
	}

	if gps_location_event.Captured_at.IsZero() {
		return fmt.Errorf("captured_at_is_required")
	}

	if gps_location_event.Latitude < -90 || gps_location_event.Latitude > 90 {
		return fmt.Errorf("latitude_out_of_range")
	}

	if gps_location_event.Longitude < -180 || gps_location_event.Longitude > 180 {
		return fmt.Errorf("longitude_out_of_range")
	}

	return nil
}
