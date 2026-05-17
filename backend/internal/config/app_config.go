package config

import (
	"fmt"
	"os"
	"strconv"
)

type App_config struct {
	Http_host                             string
	Http_port                             string
	Postgres_url                          string
	Nats_url                              string
	Osrm_base_url                         string
	Nats_gps_subject                      string
	Nats_gps_queue_group                  string
	Nats_bin_sensor_events_subject        string
	Nats_driver_collection_events_subject string
	Nats_route_blockage_events_subject    string
	Nats_route_deviation_alerts_subject   string
	Auth_jwt_signing_key                  string
	Auth_jwt_issuer                       string
	Auth_jwt_audience                     string
	Auth_access_token_ttl_minutes         int
	Auth_enable_dev_token_issue           bool
	Request_read_timeout_seconds          int
	Request_write_timeout_seconds         int
	Request_idle_timeout_seconds          int
}

func Load_app_config() (App_config, error) {
	read_timeout_seconds, read_timeout_error := read_env_int("HTTP_READ_TIMEOUT_SECONDS", 10)
	if read_timeout_error != nil {
		return App_config{}, read_timeout_error
	}

	write_timeout_seconds, write_timeout_error := read_env_int("HTTP_WRITE_TIMEOUT_SECONDS", 0)
	if write_timeout_error != nil {
		return App_config{}, write_timeout_error
	}

	idle_timeout_seconds, idle_timeout_error := read_env_int("HTTP_IDLE_TIMEOUT_SECONDS", 30)
	if idle_timeout_error != nil {
		return App_config{}, idle_timeout_error
	}

	auth_access_token_ttl_minutes, auth_access_token_ttl_minutes_error := read_env_int("AUTH_ACCESS_TOKEN_TTL_MINUTES", 120)
	if auth_access_token_ttl_minutes_error != nil {
		return App_config{}, auth_access_token_ttl_minutes_error
	}

	auth_enable_dev_token_issue, auth_enable_dev_token_issue_error := read_env_bool("AUTH_ENABLE_DEV_TOKEN_ISSUE", true)
	if auth_enable_dev_token_issue_error != nil {
		return App_config{}, auth_enable_dev_token_issue_error
	}

	application_config := App_config{
		Http_host:                             read_env_string("HTTP_HOST", "0.0.0.0"),
		Http_port:                             read_env_string("HTTP_PORT", "8080"),
		Postgres_url:                          read_env_string("POSTGRES_URL", "postgresql://ecochitas_user:ecochitas_password@localhost:5432/ecochitas_db?sslmode=disable"),
		Nats_url:                              read_env_string("NATS_URL", "nats://localhost:4222"),
		Osrm_base_url:                         read_env_string("OSRM_BASE_URL", "https://router.project-osrm.org"),
		Nats_gps_subject:                      read_env_string("NATS_GPS_SUBJECT", "gps.trucks.location"),
		Nats_gps_queue_group:                  read_env_string("NATS_GPS_QUEUE_GROUP", "ecochitas_gps_workers"),
		Nats_bin_sensor_events_subject:        read_env_string("NATS_BIN_SENSOR_EVENTS_SUBJECT", "operations.bins.sensor_events"),
		Nats_driver_collection_events_subject: read_env_string("NATS_DRIVER_COLLECTION_EVENTS_SUBJECT", "operations.driver.collection_events"),
		Nats_route_blockage_events_subject:    read_env_string("NATS_ROUTE_BLOCKAGE_EVENTS_SUBJECT", "operations.routes.blockage_events"),
		Nats_route_deviation_alerts_subject:   read_env_string("NATS_ROUTE_DEVIATION_ALERTS_SUBJECT", "operations.routes.deviation_alerts"),
		Auth_jwt_signing_key:                  read_env_string("AUTH_JWT_SIGNING_KEY", "dev_only_change_me"),
		Auth_jwt_issuer:                       read_env_string("AUTH_JWT_ISSUER", "ecochitas_backend"),
		Auth_jwt_audience:                     read_env_string("AUTH_JWT_AUDIENCE", "ecochitas_api"),
		Auth_access_token_ttl_minutes:         auth_access_token_ttl_minutes,
		Auth_enable_dev_token_issue:           auth_enable_dev_token_issue,
		Request_read_timeout_seconds:          read_timeout_seconds,
		Request_write_timeout_seconds:         write_timeout_seconds,
		Request_idle_timeout_seconds:          idle_timeout_seconds,
	}

	return application_config, nil
}

func read_env_string(environment_variable_name string, default_value string) string {
	environment_variable_value, has_environment_variable := os.LookupEnv(environment_variable_name)
	if !has_environment_variable || environment_variable_value == "" {
		return default_value
	}

	return environment_variable_value
}

func read_env_int(environment_variable_name string, default_value int) (int, error) {
	environment_variable_value, has_environment_variable := os.LookupEnv(environment_variable_name)
	if !has_environment_variable || environment_variable_value == "" {
		return default_value, nil
	}

	parsed_value, parse_error := strconv.Atoi(environment_variable_value)
	if parse_error != nil {
		return 0, fmt.Errorf("invalid_integer_environment_variable %s: %w", environment_variable_name, parse_error)
	}

	return parsed_value, nil
}

func read_env_bool(environment_variable_name string, default_value bool) (bool, error) {
	environment_variable_value, has_environment_variable := os.LookupEnv(environment_variable_name)
	if !has_environment_variable || environment_variable_value == "" {
		return default_value, nil
	}

	parsed_value, parse_error := strconv.ParseBool(environment_variable_value)
	if parse_error != nil {
		return false, fmt.Errorf("invalid_boolean_environment_variable %s: %w", environment_variable_name, parse_error)
	}

	return parsed_value, nil
}
