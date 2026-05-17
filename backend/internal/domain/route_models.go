package domain

import "time"

type Collection_route_list_query struct {
	Zone_name_filter          string
	Collection_weekday_filter int
	Has_is_active_filter      bool
	Is_active_filter_value    bool
}

type Collection_route_create_command struct {
	Route_code                    string
	Route_name                    string
	Zone_name                     string
	Collection_weekday            int
	Is_active                     bool
	Authenticated_user_identifier string
}

type Collection_route_update_command struct {
	Route_identifier              string
	Route_code                    *string
	Route_name                    *string
	Zone_name                     *string
	Collection_weekday            *int
	Is_active                     *bool
	Authenticated_user_identifier string
}

type Route_stop_sync_item struct {
	Bin_identifier string
	Stop_order     int
	Planned_time   string
}

type Route_stop_sync_command struct {
	Route_identifier              string
	Stop_item_list                []Route_stop_sync_item
	Authenticated_user_identifier string
}

type Route_path_coordinate struct {
	Stop_order     int     `json:"stop_order"`
	Bin_identifier string  `json:"bin_identifier"`
	Bin_code       string  `json:"bin_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}

type Route_road_path_coordinate [2]float64

type Collection_route_view struct {
	Route_identifier      string                       `json:"route_identifier"`
	Route_code            string                       `json:"route_code"`
	Route_name            string                       `json:"route_name"`
	Zone_name             string                       `json:"zone_name"`
	Collection_weekday    int                          `json:"collection_weekday"`
	Is_active             bool                         `json:"is_active"`
	Stop_total            int                          `json:"stop_total"`
	Path_coordinates      []Route_path_coordinate      `json:"path_coordinates"`
	Road_path_coordinates []Route_road_path_coordinate `json:"road_path_coordinates"`
	Routing_status        string                       `json:"routing_status"`
	Routing_provider      string                       `json:"routing_provider"`
	Routed_at             *time.Time                   `json:"routed_at,omitempty"`
}

type Route_stop_view struct {
	Route_stop_identifier string  `json:"route_stop_identifier"`
	Route_identifier      string  `json:"route_identifier"`
	Bin_identifier        string  `json:"bin_identifier"`
	Bin_code              string  `json:"bin_code"`
	Zone_name             string  `json:"zone_name"`
	Stop_order            int     `json:"stop_order"`
	Planned_time          string  `json:"planned_time,omitempty"`
	Latitude              float64 `json:"latitude"`
	Longitude             float64 `json:"longitude"`
}

type Collection_route_revision_list_query struct {
	Route_identifier string
	Limit            int
}

type Collection_route_revision_record struct {
	Revision_identifier        string         `json:"revision_identifier"`
	Route_identifier           string         `json:"route_identifier"`
	Revision_number            int            `json:"revision_number"`
	Change_type                string         `json:"change_type"`
	Changed_by_user_identifier string         `json:"changed_by_user_identifier,omitempty"`
	Change_payload             map[string]any `json:"change_payload"`
	Created_at                 time.Time      `json:"created_at"`
}

type Collection_route_revision_list_result struct {
	Items []Collection_route_revision_record `json:"items"`
	Total int                                `json:"total"`
}

type Truck_route_deviation_query struct {
	Truck_identifier           string
	Route_identifier           string
	Deviation_threshold_meters float64
}

type Truck_route_deviation_view struct {
	Truck_identifier              string    `json:"truck_identifier"`
	Route_identifier              string    `json:"route_identifier"`
	Captured_at                   time.Time `json:"captured_at"`
	Truck_latitude                float64   `json:"truck_latitude"`
	Truck_longitude               float64   `json:"truck_longitude"`
	Nearest_route_stop_identifier string    `json:"nearest_route_stop_identifier"`
	Nearest_stop_order            int       `json:"nearest_stop_order"`
	Nearest_bin_identifier        string    `json:"nearest_bin_identifier"`
	Nearest_bin_code              string    `json:"nearest_bin_code"`
	Nearest_bin_latitude          float64   `json:"nearest_bin_latitude"`
	Nearest_bin_longitude         float64   `json:"nearest_bin_longitude"`
	Distance_to_route_meters      float64   `json:"distance_to_route_meters"`
	Deviation_threshold_meters    float64   `json:"deviation_threshold_meters"`
	Is_off_route                  bool      `json:"is_off_route"`
}

type Truck_route_assignment_create_command struct {
	Truck_identifier            string
	Route_identifier            string
	Assigned_by_user_identifier string
	Assignment_notes            string
}

type Truck_route_assignment_list_query struct {
	Truck_identifier string
	Route_identifier string
	Is_active_filter *bool
	Limit            int
}

type Truck_route_assignment_view struct {
	Assignment_identifier       string     `json:"assignment_identifier"`
	Truck_identifier            string     `json:"truck_identifier"`
	Route_identifier            string     `json:"route_identifier"`
	Route_code                  string     `json:"route_code"`
	Route_name                  string     `json:"route_name"`
	Zone_name                   string     `json:"zone_name"`
	Collection_weekday          int        `json:"collection_weekday"`
	Is_active                   bool       `json:"is_active"`
	Assigned_by_user_identifier string     `json:"assigned_by_user_identifier,omitempty"`
	Assignment_notes            string     `json:"assignment_notes,omitempty"`
	Assigned_at                 time.Time  `json:"assigned_at"`
	Unassigned_at               *time.Time `json:"unassigned_at,omitempty"`
}

type Truck_route_assignment_list_result struct {
	Items []Truck_route_assignment_view `json:"items"`
	Total int                           `json:"total"`
}

type Route_deviation_alert_create_command struct {
	Truck_identifier             string
	Route_identifier             string
	Deviation_threshold_meters   float64
	Triggered_by_user_identifier string
	Alert_notes                  string
}

type Route_deviation_alert_list_query struct {
	Truck_identifier string
	Route_identifier string
	Alert_status     string
	Limit            int
}

type Route_deviation_alert_record struct {
	Alert_identifier             string         `json:"alert_identifier"`
	Truck_identifier             string         `json:"truck_identifier"`
	Route_identifier             string         `json:"route_identifier"`
	Route_stop_identifier        string         `json:"route_stop_identifier,omitempty"`
	Distance_to_route_meters     float64        `json:"distance_to_route_meters"`
	Deviation_threshold_meters   float64        `json:"deviation_threshold_meters"`
	Severity_level               string         `json:"severity_level"`
	Alert_status                 string         `json:"alert_status"`
	Alert_notes                  string         `json:"alert_notes,omitempty"`
	Triggered_by_user_identifier string         `json:"triggered_by_user_identifier,omitempty"`
	Detected_at                  time.Time      `json:"detected_at"`
	Resolved_at                  *time.Time     `json:"resolved_at,omitempty"`
	Metadata_payload             map[string]any `json:"metadata_payload"`
	Created_at                   time.Time      `json:"created_at"`
}

type Route_deviation_alert_list_result struct {
	Items []Route_deviation_alert_record `json:"items"`
	Total int                            `json:"total"`
}
