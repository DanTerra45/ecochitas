package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecochitas/internal/domain"

	"github.com/nats-io/nats.go"
)

type waypoint struct {
	latitude  float64
	longitude float64
}

// Demo route waypoints (Cochabamba Centro) — match BIN-DEMO-001/002/003.
var default_route_waypoints = []waypoint{
	{latitude: -17.3935, longitude: -66.1571},
	{latitude: -17.3917, longitude: -66.1548},
	{latitude: -17.3899, longitude: -66.1527},
}

func main() {
	nats_url_flag := flag.String("nats-url", read_env_string("NATS_URL", "nats://localhost:4222"), "NATS server URL")
	gps_subject_flag := flag.String("subject", read_env_string("NATS_GPS_SUBJECT", "gps.trucks.location"), "NATS subject to publish GPS events to")
	truck_identifier_flag := flag.String("truck", read_env_string("GPS_TRUCK_IDENTIFIER", "TRUCK-001"), "Truck identifier")
	tick_interval_flag := flag.Duration("interval", 1*time.Second, "Interval between published positions")
	steps_per_segment_flag := flag.Int("steps", 20, "Interpolation steps between consecutive waypoints")
	loop_route_flag := flag.Bool("loop", true, "Loop the route forever (otherwise publish one pass and exit)")
	flag.Parse()

	nats_connection, connect_nats_error := nats.Connect(*nats_url_flag)
	if connect_nats_error != nil {
		slog.Error("failed_to_connect_nats", "error", connect_nats_error)
		os.Exit(1)
	}
	defer nats_connection.Close()

	slog.Info(
		"gps_simulator_started",
		"truck_identifier", *truck_identifier_flag,
		"subject", *gps_subject_flag,
		"interval", tick_interval_flag.String(),
		"steps_per_segment", *steps_per_segment_flag,
		"loop", *loop_route_flag,
	)

	stop_signal_channel := make(chan os.Signal, 1)
	signal.Notify(stop_signal_channel, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(*tick_interval_flag)
	defer ticker.Stop()

	path := build_interpolated_path(default_route_waypoints, *steps_per_segment_flag)
	if len(path) == 0 {
		slog.Error("empty_simulation_path")
		os.Exit(1)
	}

	current_index := 0
	previous_point := path[0]

	for {
		select {
		case <-stop_signal_channel:
			slog.Info("gps_simulator_stopping")
			return
		case <-ticker.C:
			current_point := path[current_index]
			heading_degrees := compute_heading_degrees(previous_point, current_point)
			speed_kmh := 22.0 + 6.0*math.Sin(float64(current_index)/4.0)

			gps_location_event := domain.Gps_location_event{
				Truck_identifier: *truck_identifier_flag,
				Latitude:         current_point.latitude,
				Longitude:        current_point.longitude,
				Speed_kmh:        &speed_kmh,
				Heading_degrees:  &heading_degrees,
				Captured_at:      time.Now().UTC(),
			}

			serialized_payload, marshal_error := json.Marshal(gps_location_event)
			if marshal_error != nil {
				slog.Error("failed_to_marshal_payload", "error", marshal_error)
				continue
			}

			if publish_error := nats_connection.Publish(*gps_subject_flag, serialized_payload); publish_error != nil {
				slog.Error("failed_to_publish_gps_event", "error", publish_error)
				continue
			}

			slog.Info(
				"gps_event_published",
				"truck_identifier", gps_location_event.Truck_identifier,
				"latitude", gps_location_event.Latitude,
				"longitude", gps_location_event.Longitude,
				"step", current_index+1,
				"total_steps", len(path),
			)

			previous_point = current_point
			current_index++
			if current_index >= len(path) {
				if !*loop_route_flag {
					slog.Info("gps_simulator_finished_one_pass")
					return
				}
				current_index = 0
			}
		}
	}
}

func build_interpolated_path(route_waypoints []waypoint, steps_per_segment int) []waypoint {
	if len(route_waypoints) < 2 || steps_per_segment < 1 {
		return route_waypoints
	}

	// Match the frontend's "manhattan" mock: between every pair of stops
	// insert an L-shaped elbow, alternating which axis we travel first.
	manhattan_waypoints := []waypoint{route_waypoints[0]}
	for segment_index := 0; segment_index < len(route_waypoints)-1; segment_index++ {
		segment_start := route_waypoints[segment_index]
		segment_end := route_waypoints[segment_index+1]
		var elbow_point waypoint
		if segment_index%2 == 0 {
			elbow_point = waypoint{latitude: segment_start.latitude, longitude: segment_end.longitude}
		} else {
			elbow_point = waypoint{latitude: segment_end.latitude, longitude: segment_start.longitude}
		}
		manhattan_waypoints = append(manhattan_waypoints, elbow_point, segment_end)
	}

	interpolated_path := make([]waypoint, 0, len(manhattan_waypoints)*steps_per_segment)
	for waypoint_index := 0; waypoint_index < len(manhattan_waypoints)-1; waypoint_index++ {
		segment_start := manhattan_waypoints[waypoint_index]
		segment_end := manhattan_waypoints[waypoint_index+1]
		for step_index := 0; step_index < steps_per_segment; step_index++ {
			progress_fraction := float64(step_index) / float64(steps_per_segment)
			interpolated_path = append(interpolated_path, waypoint{
				latitude:  segment_start.latitude + (segment_end.latitude-segment_start.latitude)*progress_fraction,
				longitude: segment_start.longitude + (segment_end.longitude-segment_start.longitude)*progress_fraction,
			})
		}
	}
	interpolated_path = append(interpolated_path, manhattan_waypoints[len(manhattan_waypoints)-1])
	return interpolated_path
}

func compute_heading_degrees(from_point waypoint, to_point waypoint) float64 {
	delta_longitude_radians := (to_point.longitude - from_point.longitude) * math.Pi / 180
	from_latitude_radians := from_point.latitude * math.Pi / 180
	to_latitude_radians := to_point.latitude * math.Pi / 180

	y_component := math.Sin(delta_longitude_radians) * math.Cos(to_latitude_radians)
	x_component := math.Cos(from_latitude_radians)*math.Sin(to_latitude_radians) -
		math.Sin(from_latitude_radians)*math.Cos(to_latitude_radians)*math.Cos(delta_longitude_radians)

	heading_radians := math.Atan2(y_component, x_component)
	heading_degrees := heading_radians * 180 / math.Pi
	if heading_degrees < 0 {
		heading_degrees += 360
	}
	return heading_degrees
}

func read_env_string(environment_variable_name string, default_value string) string {
	environment_variable_value := os.Getenv(environment_variable_name)
	if environment_variable_value == "" {
		return default_value
	}
	return environment_variable_value
}
