package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"ecochitas/internal/domain"

	"github.com/nats-io/nats.go"
)

func main() {
	nats_url := read_env_string("NATS_URL", "nats://localhost:4222")
	gps_subject := read_env_string("NATS_GPS_SUBJECT", "gps.trucks.location")

	nats_connection, connect_nats_error := nats.Connect(nats_url)
	if connect_nats_error != nil {
		slog.Error("failed_to_connect_nats", "error", connect_nats_error)
		os.Exit(1)
	}
	defer nats_connection.Close()

	sample_speed_kmh := 28.4
	sample_heading_degrees := 142.0
	gps_location_event := domain.Gps_location_event{
		Truck_identifier: "TRUCK-001",
		Latitude:         -17.3935,
		Longitude:        -66.1570,
		Speed_kmh:        &sample_speed_kmh,
		Heading_degrees:  &sample_heading_degrees,
		Captured_at:      time.Now().UTC(),
	}

	serialized_payload, marshal_payload_error := json.Marshal(gps_location_event)
	if marshal_payload_error != nil {
		slog.Error("failed_to_marshal_payload", "error", marshal_payload_error)
		os.Exit(1)
	}

	publish_event_error := nats_connection.Publish(gps_subject, serialized_payload)
	if publish_event_error != nil {
		slog.Error("failed_to_publish_gps_event", "error", publish_event_error)
		os.Exit(1)
	}

	slog.Info("gps_event_published", "subject", gps_subject, "truck_identifier", gps_location_event.Truck_identifier)
}

func read_env_string(environment_variable_name string, default_value string) string {
	environment_variable_value := os.Getenv(environment_variable_name)
	if environment_variable_value == "" {
		return default_value
	}

	return environment_variable_value
}
