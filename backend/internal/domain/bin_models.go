package domain

import "time"

type Bin_status struct {
	Bin_identifier   string     `json:"bin_identifier"`
	Bin_code         string     `json:"bin_code"`
	Zone_name        string     `json:"zone_name"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	Fill_percentage  int        `json:"fill_percentage"`
	Sensor_status    string     `json:"sensor_status"`
	Bin_status_label string     `json:"bin_status"`
	Last_emptied_at  *time.Time `json:"last_emptied_at,omitempty"`
}
