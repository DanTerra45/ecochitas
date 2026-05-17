package domain

import "time"

type Zone_recycling_container_view struct {
	Container_identifier     string  `json:"container_identifier"`
	Container_code           string  `json:"container_code"`
	Zone_identifier          string  `json:"zone_identifier"`
	Zone_code                string  `json:"zone_code"`
	Zone_name                string  `json:"zone_name"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
	Assigned_household_total int     `json:"assigned_household_total"`
	Eligible_household_total int     `json:"eligible_household_total"`
	Contaminated_cycle_total int     `json:"contaminated_cycle_total"`
	Completed_cycle_total    int     `json:"completed_cycle_total"`
}

type Recycling_cycle_start_command struct {
	Container_identifier      string    `json:"container_identifier"`
	Scheduled_collection_date time.Time `json:"scheduled_collection_date"`
	Collection_operator_name  string    `json:"collection_operator_name,omitempty"`
}

type Recycling_collection_cycle struct {
	Cycle_identifier               string     `json:"cycle_identifier"`
	Container_identifier           string     `json:"container_identifier"`
	Container_code                 string     `json:"container_code"`
	Zone_identifier                string     `json:"zone_identifier"`
	Zone_code                      string     `json:"zone_code"`
	Zone_name                      string     `json:"zone_name"`
	Scheduled_collection_date      string     `json:"scheduled_collection_date"`
	Cycle_status                   string     `json:"cycle_status"`
	Collection_operator_name       string     `json:"collection_operator_name,omitempty"`
	Raw_points_total               float64    `json:"raw_points_total"`
	Contamination_level            string     `json:"contamination_level"`
	Contamination_discount_percent float64    `json:"contamination_discount_percent"`
	Discount_points_total          float64    `json:"discount_points_total"`
	Final_points_total             float64    `json:"final_points_total"`
	Contamination_notes            string     `json:"contamination_notes,omitempty"`
	Closed_at                      *time.Time `json:"closed_at,omitempty"`
}

type Recycling_evidence_submission_command struct {
	Cycle_identifier     string    `json:"cycle_identifier"`
	Household_identifier string    `json:"household_identifier"`
	Evidence_photo_url   string    `json:"evidence_photo_url"`
	Evidence_captured_at time.Time `json:"evidence_captured_at"`
	Evidence_latitude    float64   `json:"evidence_latitude"`
	Evidence_longitude   float64   `json:"evidence_longitude"`
	Validation_status    string    `json:"validation_status,omitempty"`
	Rejection_reason     string    `json:"rejection_reason,omitempty"`
}

type Recycling_evidence_submission struct {
	Submission_identifier string    `json:"submission_identifier"`
	Cycle_identifier      string    `json:"cycle_identifier"`
	Household_identifier  string    `json:"household_identifier"`
	Validation_status     string    `json:"validation_status"`
	Created_at            time.Time `json:"created_at"`
	Updated_at            time.Time `json:"updated_at"`
}

type Recycling_cycle_close_command struct {
	Cycle_identifier                  string   `json:"cycle_identifier"`
	Raw_points_total                  float64  `json:"raw_points_total"`
	Contamination_level               string   `json:"contamination_level"`
	Contamination_discount_percentage *float64 `json:"contamination_discount_percentage,omitempty"`
	Contamination_notes               string   `json:"contamination_notes,omitempty"`
	Collection_operator_name          string   `json:"collection_operator_name,omitempty"`
}

type Recycling_cycle_household_points struct {
	Household_identifier string  `json:"household_identifier"`
	Household_code       string  `json:"household_code"`
	Awarded_points       float64 `json:"awarded_points"`
}

type Recycling_cycle_summary struct {
	Cycle_identifier                  string                             `json:"cycle_identifier"`
	Container_identifier              string                             `json:"container_identifier"`
	Container_code                    string                             `json:"container_code"`
	Zone_identifier                   string                             `json:"zone_identifier"`
	Zone_code                         string                             `json:"zone_code"`
	Zone_name                         string                             `json:"zone_name"`
	Cycle_status                      string                             `json:"cycle_status"`
	Scheduled_collection_date         string                             `json:"scheduled_collection_date"`
	Collection_operator_name          string                             `json:"collection_operator_name,omitempty"`
	Eligible_household_total          int                                `json:"eligible_household_total"`
	Raw_points_total                  float64                            `json:"raw_points_total"`
	Contamination_level               string                             `json:"contamination_level"`
	Contamination_discount_percentage float64                            `json:"contamination_discount_percentage"`
	Discount_points_total             float64                            `json:"discount_points_total"`
	Final_points_total                float64                            `json:"final_points_total"`
	Contamination_notes               string                             `json:"contamination_notes,omitempty"`
	Closed_at                         *time.Time                         `json:"closed_at,omitempty"`
	Household_points                  []Recycling_cycle_household_points `json:"household_points"`
}
