package domain

import "time"

type Driver_collection_event_create_command struct {
	Route_stop_identifier         string
	Bin_identifier                string
	Authenticated_user_identifier string
	Authenticated_role_name       string
	Authenticated_full_name       string
	Action_type                   string
	Evidence_photo_url            string
	Action_notes                  string
	Action_at                     time.Time
}

type Driver_collection_event_record struct {
	Event_identifier       string     `json:"event_identifier"`
	Bin_identifier         string     `json:"bin_identifier"`
	Route_stop_identifier  string     `json:"route_stop_identifier,omitempty"`
	Driver_user_identifier string     `json:"driver_user_identifier"`
	Action_type            string     `json:"action_type"`
	Evidence_photo_url     string     `json:"evidence_photo_url,omitempty"`
	Action_notes           string     `json:"action_notes,omitempty"`
	Action_at              time.Time  `json:"action_at"`
	Bin_fill_percentage    int        `json:"bin_fill_percentage"`
	Bin_status             string     `json:"bin_status"`
	Bin_sensor_status      string     `json:"bin_sensor_status"`
	Bin_last_emptied_at    *time.Time `json:"bin_last_emptied_at,omitempty"`
}

type Bin_sensor_event_ingestion_command struct {
	Bin_identifier    string
	Source_identifier string
	Fill_percentage   int
	Sensor_status     string
	Measured_at       time.Time
}

type Bin_sensor_event_record struct {
	Sensor_event_identifier string     `json:"sensor_event_identifier"`
	Bin_identifier          string     `json:"bin_identifier"`
	Source_identifier       string     `json:"source_identifier,omitempty"`
	Fill_percentage         int        `json:"fill_percentage"`
	Sensor_status           string     `json:"sensor_status"`
	Bin_status              string     `json:"bin_status"`
	Measured_at             time.Time  `json:"measured_at"`
	Recorded_at             time.Time  `json:"recorded_at"`
	Bin_last_emptied_at     *time.Time `json:"bin_last_emptied_at,omitempty"`
}

type Route_blockage_report_create_command struct {
	Route_identifier              string
	Route_stop_identifier         string
	Bin_identifier                string
	Authenticated_user_identifier string
	Blockage_reason               string
	Evidence_photo_url            string
	Severity_level                string
	Reported_at                   time.Time
}

type Route_blockage_report_list_query struct {
	Route_identifier      string
	Route_stop_identifier string
	Bin_identifier        string
	Status_filter         string
	Limit                 int
}

type Route_blockage_report_status_update_command struct {
	Blockage_identifier           string
	Authenticated_user_identifier string
	Status                        string
	Resolution_notes              string
	Resolved_at                   time.Time
}

type Route_blockage_report_record struct {
	Blockage_identifier         string     `json:"blockage_identifier"`
	Route_identifier            string     `json:"route_identifier,omitempty"`
	Route_stop_identifier       string     `json:"route_stop_identifier,omitempty"`
	Bin_identifier              string     `json:"bin_identifier,omitempty"`
	Reported_by_user_identifier string     `json:"reported_by_user_identifier"`
	Blockage_reason             string     `json:"blockage_reason"`
	Evidence_photo_url          string     `json:"evidence_photo_url,omitempty"`
	Severity_level              string     `json:"severity_level"`
	Status                      string     `json:"status"`
	Reported_at                 time.Time  `json:"reported_at"`
	Resolved_at                 *time.Time `json:"resolved_at,omitempty"`
	Resolution_notes            string     `json:"resolution_notes,omitempty"`
	Created_at                  time.Time  `json:"created_at"`
	Updated_at                  time.Time  `json:"updated_at"`
}

type Route_blockage_report_list_result struct {
	Items []Route_blockage_report_record `json:"items"`
	Total int                            `json:"total"`
}
