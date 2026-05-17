package domain

import "time"

type Gps_location_event struct {
	Truck_identifier string    `json:"truck_identifier"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Speed_kmh        *float64  `json:"speed_kmh,omitempty"`
	Heading_degrees  *float64  `json:"heading_degrees,omitempty"`
	Captured_at      time.Time `json:"captured_at"`
}

type Truck_latest_position struct {
	Truck_identifier string    `json:"truck_identifier"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Speed_kmh        *float64  `json:"speed_kmh,omitempty"`
	Heading_degrees  *float64  `json:"heading_degrees,omitempty"`
	Captured_at      time.Time `json:"captured_at"`
	Received_at      time.Time `json:"received_at"`
}
