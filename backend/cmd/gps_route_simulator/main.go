package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ecochitas/internal/domain"

	"github.com/nats-io/nats.go"
)

type Route_simulator_config struct {
	Nats_url            string
	Gps_subject         string
	Truck_identifier    string
	Center_latitude     float64
	Center_longitude    float64
	Path_radius_degrees float64
	Publish_interval    time.Duration
	Rotation_steps      int
	Speed_kmh           float64
}

func main() {
	route_simulator_config := load_route_simulator_config()

	nats_connection, connect_nats_error := nats.Connect(route_simulator_config.Nats_url)
	if connect_nats_error != nil {
		slog.Error("failed_to_connect_nats", "error", connect_nats_error)
		os.Exit(1)
	}
	defer nats_connection.Close()

	slog.Info(
		"gps_route_simulator_started",
		"truck_identifier",
		route_simulator_config.Truck_identifier,
		"publish_interval",
		route_simulator_config.Publish_interval,
		"gps_subject",
		route_simulator_config.Gps_subject,
	)

	application_context, cancel_signal_context := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel_signal_context()

	publish_ticker := time.NewTicker(route_simulator_config.Publish_interval)
	defer publish_ticker.Stop()

	route_step_index := 0
	for {
		select {
		case <-application_context.Done():
			slog.Info("gps_route_simulator_stopped")
			return
		case captured_time := <-publish_ticker.C:
			gps_location_event := build_simulated_gps_event(
				route_simulator_config,
				route_step_index,
				captured_time.UTC(),
			)

			publish_error := publish_gps_event(nats_connection, route_simulator_config.Gps_subject, gps_location_event)
			if publish_error != nil {
				slog.Error("failed_to_publish_simulated_gps_event", "error", publish_error)
				continue
			}

			slog.Info(
				"simulated_gps_event_published",
				"truck_identifier",
				gps_location_event.Truck_identifier,
				"latitude",
				gps_location_event.Latitude,
				"longitude",
				gps_location_event.Longitude,
			)

			route_step_index += 1
		}
	}
}

func load_route_simulator_config() Route_simulator_config {
	truck_identifier_flag := flag.String("truck_identifier", "TRUCK-001", "truck identifier")
	center_latitude_flag := flag.Float64("center_latitude", -17.3935, "route center latitude")
	center_longitude_flag := flag.Float64("center_longitude", -66.1570, "route center longitude")
	path_radius_degrees_flag := flag.Float64("path_radius_degrees", 0.0025, "route radius in degrees")
	publish_interval_flag := flag.Duration("publish_interval", time.Second, "publish interval, e.g. 1s, 500ms")
	rotation_steps_flag := flag.Int("rotation_steps", 120, "steps for full circular rotation")
	speed_kmh_flag := flag.Float64("speed_kmh", 28.4, "simulated speed km/h")
	flag.Parse()

	nats_url := read_env_string("NATS_URL", "nats://localhost:4222")
	gps_subject := read_env_string("NATS_GPS_SUBJECT", "gps.trucks.location")
	truck_identifier := strings.TrimSpace(*truck_identifier_flag)
	if truck_identifier == "" {
		truck_identifier = "TRUCK-001"
	}

	rotation_steps := *rotation_steps_flag
	if rotation_steps < 12 {
		rotation_steps = 12
	}

	publish_interval := *publish_interval_flag
	if publish_interval < 200*time.Millisecond {
		publish_interval = 200 * time.Millisecond
	}

	return Route_simulator_config{
		Nats_url:            nats_url,
		Gps_subject:         gps_subject,
		Truck_identifier:    truck_identifier,
		Center_latitude:     *center_latitude_flag,
		Center_longitude:    *center_longitude_flag,
		Path_radius_degrees: *path_radius_degrees_flag,
		Publish_interval:    publish_interval,
		Rotation_steps:      rotation_steps,
		Speed_kmh:           *speed_kmh_flag,
	}
}

func build_simulated_gps_event(
	route_simulator_config Route_simulator_config,
	route_step_index int,
	captured_time time.Time,
) domain.Gps_location_event {
	rotation_position := route_step_index % route_simulator_config.Rotation_steps
	rotation_progress := float64(rotation_position) / float64(route_simulator_config.Rotation_steps)
	angle_radians := 2 * math.Pi * rotation_progress

	latitude_offset := route_simulator_config.Path_radius_degrees * math.Sin(angle_radians)
	longitude_offset := route_simulator_config.Path_radius_degrees * math.Cos(angle_radians)

	simulated_latitude := route_simulator_config.Center_latitude + latitude_offset
	simulated_longitude := route_simulator_config.Center_longitude + longitude_offset

	heading_degrees := normalize_heading_degrees(angle_radians * (180 / math.Pi))
	speed_kmh := route_simulator_config.Speed_kmh

	return domain.Gps_location_event{
		Truck_identifier: route_simulator_config.Truck_identifier,
		Latitude:         simulated_latitude,
		Longitude:        simulated_longitude,
		Speed_kmh:        &speed_kmh,
		Heading_degrees:  &heading_degrees,
		Captured_at:      captured_time,
	}
}

func publish_gps_event(
	nats_connection *nats.Conn,
	gps_subject string,
	gps_location_event domain.Gps_location_event,
) error {
	serialized_payload, marshal_payload_error := json.Marshal(gps_location_event)
	if marshal_payload_error != nil {
		return marshal_payload_error
	}

	return nats_connection.Publish(gps_subject, serialized_payload)
}

func normalize_heading_degrees(raw_heading_degrees float64) float64 {
	normalized_heading := math.Mod(raw_heading_degrees, 360)
	if normalized_heading < 0 {
		normalized_heading += 360
	}

	return normalized_heading
}

func read_env_string(environment_variable_name string, default_value string) string {
	environment_variable_value := os.Getenv(environment_variable_name)
	if environment_variable_value == "" {
		return default_value
	}

	return environment_variable_value
}
