package config

import (
	"fmt"
	"os"
	"strconv"
)

type App_config struct {
	Http_host                     string
	Http_port                     string
	Postgres_url                  string
	Nats_url                      string
	Nats_gps_subject              string
	Nats_gps_queue_group          string
	Request_read_timeout_seconds  int
	Request_write_timeout_seconds int
	Request_idle_timeout_seconds  int
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

	application_config := App_config{
		Http_host:                     read_env_string("HTTP_HOST", "0.0.0.0"),
		Http_port:                     read_env_string("HTTP_PORT", "8080"),
		Postgres_url:                  read_env_string("POSTGRES_URL", "postgresql://ecochitas_user:ecochitas_password@localhost:5432/ecochitas_db?sslmode=disable"),
		Nats_url:                      read_env_string("NATS_URL", "nats://localhost:4222"),
		Nats_gps_subject:              read_env_string("NATS_GPS_SUBJECT", "gps.trucks.location"),
		Nats_gps_queue_group:          read_env_string("NATS_GPS_QUEUE_GROUP", "ecochitas_gps_workers"),
		Request_read_timeout_seconds:  read_timeout_seconds,
		Request_write_timeout_seconds: write_timeout_seconds,
		Request_idle_timeout_seconds:  idle_timeout_seconds,
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
