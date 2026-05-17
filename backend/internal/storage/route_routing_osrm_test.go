package storage

import (
	"testing"

	"ecochitas/internal/domain"
)

func Test_build_route_stops_hash_is_deterministic_for_same_ordered_coordinates(t *testing.T) {
	first_hash := build_route_stops_hash([]domain.Route_path_coordinate{
		{
			Stop_order: 1,
			Latitude:   -17.3935000,
			Longitude:  -66.1570000,
		},
		{
			Stop_order: 2,
			Latitude:   -17.3901000,
			Longitude:  -66.1602000,
		},
	})
	second_hash := build_route_stops_hash([]domain.Route_path_coordinate{
		{
			Stop_order: 1,
			Latitude:   -17.3935000,
			Longitude:  -66.1570000,
		},
		{
			Stop_order: 2,
			Latitude:   -17.3901000,
			Longitude:  -66.1602000,
		},
	})

	if first_hash == "" || second_hash == "" {
		t.Fatalf("expected non-empty hash values")
	}
	if first_hash != second_hash {
		t.Fatalf("expected deterministic hash for identical ordered stops")
	}
}

func Test_build_route_stops_hash_changes_when_stop_order_changes(t *testing.T) {
	forward_hash := build_route_stops_hash([]domain.Route_path_coordinate{
		{
			Stop_order: 1,
			Latitude:   -17.3935000,
			Longitude:  -66.1570000,
		},
		{
			Stop_order: 2,
			Latitude:   -17.3901000,
			Longitude:  -66.1602000,
		},
	})
	reversed_hash := build_route_stops_hash([]domain.Route_path_coordinate{
		{
			Stop_order: 2,
			Latitude:   -17.3901000,
			Longitude:  -66.1602000,
		},
		{
			Stop_order: 1,
			Latitude:   -17.3935000,
			Longitude:  -66.1570000,
		},
	})

	if forward_hash == reversed_hash {
		t.Fatalf("expected different hashes when ordered stops change")
	}
}

func Test_build_straight_line_road_path_coordinates_maps_stop_coordinates(t *testing.T) {
	straight_line_coordinates := build_straight_line_road_path_coordinates([]domain.Route_path_coordinate{
		{
			Stop_order: 1,
			Latitude:   -17.3935000,
			Longitude:  -66.1570000,
		},
		{
			Stop_order: 2,
			Latitude:   -17.3901000,
			Longitude:  -66.1602000,
		},
	})

	if len(straight_line_coordinates) != 2 {
		t.Fatalf("expected 2 straight-line coordinates, got %d", len(straight_line_coordinates))
	}
	if straight_line_coordinates[0][0] != -66.1570000 || straight_line_coordinates[0][1] != -17.3935000 {
		t.Fatalf("unexpected first coordinate mapping")
	}
	if straight_line_coordinates[1][0] != -66.1602000 || straight_line_coordinates[1][1] != -17.3901000 {
		t.Fatalf("unexpected second coordinate mapping")
	}
}

func Test_normalize_osrm_base_url(t *testing.T) {
	normalized_from_empty := normalize_osrm_base_url("")
	if normalized_from_empty != default_osrm_base_url {
		t.Fatalf("expected default osrm base url when empty, got %s", normalized_from_empty)
	}

	normalized_with_trailing_slash := normalize_osrm_base_url("https://example.com///")
	if normalized_with_trailing_slash != "https://example.com" {
		t.Fatalf("expected trailing slash removal, got %s", normalized_with_trailing_slash)
	}
}
