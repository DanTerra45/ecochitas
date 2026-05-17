package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"ecochitas/internal/api"
	"ecochitas/internal/auth"
	"ecochitas/internal/config"
	"ecochitas/internal/gps"
	"ecochitas/internal/infrastructure/nats_client"
	"ecochitas/internal/infrastructure/postgres"
	"ecochitas/internal/realtime"
	"ecochitas/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type Backend_app struct {
	application_config  config.App_config
	application_logger  *slog.Logger
	postgres_pool       *pgxpool.Pool
	nats_connection     *nats.Conn
	gps_event_consumer  *gps.Gps_event_consumer
	api_server_instance *api.Api_server
}

func New_backend_app() (*Backend_app, error) {
	application_logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	application_config, load_config_error := config.Load_app_config()
	if load_config_error != nil {
		return nil, fmt.Errorf("failed_to_load_application_config: %w", load_config_error)
	}

	startup_context, cancel_startup_context := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel_startup_context()

	postgres_pool, connect_postgres_error := postgres.New_postgres_pool(startup_context, application_config.Postgres_url)
	if connect_postgres_error != nil {
		return nil, connect_postgres_error
	}

	nats_connection, connect_nats_error := nats_client.New_nats_connection(application_config.Nats_url)
	if connect_nats_error != nil {
		postgres_pool.Close()
		return nil, connect_nats_error
	}

	bin_repository := storage.New_bin_repository(postgres_pool)
	truck_position_repository := storage.New_truck_position_repository(postgres_pool)
	route_repository := storage.New_route_repository(postgres_pool, application_config.Osrm_base_url)
	recycling_zone_repository := storage.New_recycling_zone_repository(postgres_pool)
	operations_repository := storage.New_operations_repository(postgres_pool)
	truck_position_stream := realtime.New_truck_position_stream(application_logger)
	operations_event_publisher := realtime.New_operations_event_publisher(
		nats_connection,
		application_config.Nats_bin_sensor_events_subject,
		application_config.Nats_driver_collection_events_subject,
		application_config.Nats_route_blockage_events_subject,
		application_config.Nats_route_deviation_alerts_subject,
	)
	jwt_authenticator, create_jwt_authenticator_error := auth.New_jwt_authenticator(
		application_config.Auth_jwt_signing_key,
		application_config.Auth_jwt_issuer,
		application_config.Auth_jwt_audience,
		application_config.Auth_access_token_ttl_minutes,
	)
	if create_jwt_authenticator_error != nil {
		nats_connection.Close()
		postgres_pool.Close()
		return nil, fmt.Errorf("failed_to_create_jwt_authenticator: %w", create_jwt_authenticator_error)
	}

	api_handler := api.New_api_handler(
		bin_repository,
		truck_position_repository,
		route_repository,
		recycling_zone_repository,
		operations_repository,
		truck_position_stream,
		operations_event_publisher,
		jwt_authenticator,
		application_config.Auth_enable_dev_token_issue,
		application_logger,
	)
	api_server_instance := api.New_api_server(application_config, api_handler, application_logger)
	gps_event_consumer := gps.New_gps_event_consumer(
		nats_connection,
		application_config.Nats_gps_subject,
		application_config.Nats_gps_queue_group,
		truck_position_repository,
		truck_position_stream,
		application_logger,
	)

	return &Backend_app{
		application_config:  application_config,
		application_logger:  application_logger,
		postgres_pool:       postgres_pool,
		nats_connection:     nats_connection,
		gps_event_consumer:  gps_event_consumer,
		api_server_instance: api_server_instance,
	}, nil
}

func (backend_app *Backend_app) Run(application_context context.Context) error {
	start_consumer_error := backend_app.gps_event_consumer.Start()
	if start_consumer_error != nil {
		return start_consumer_error
	}

	server_error_channel := make(chan error, 1)
	go func() {
		server_error_channel <- backend_app.api_server_instance.Start()
	}()

	select {
	case <-application_context.Done():
		shutdown_context, cancel_shutdown_context := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel_shutdown_context()
		return backend_app.shutdown(shutdown_context)
	case server_error := <-server_error_channel:
		return server_error
	}
}

func (backend_app *Backend_app) shutdown(shutdown_context context.Context) error {
	shutdown_error_list := make([]error, 0)

	shutdown_http_error := backend_app.api_server_instance.Shutdown(shutdown_context)
	if shutdown_http_error != nil {
		shutdown_error_list = append(shutdown_error_list, fmt.Errorf("failed_to_shutdown_http_server: %w", shutdown_http_error))
	}

	stop_gps_consumer_error := backend_app.gps_event_consumer.Stop()
	if stop_gps_consumer_error != nil {
		shutdown_error_list = append(shutdown_error_list, fmt.Errorf("failed_to_stop_gps_consumer: %w", stop_gps_consumer_error))
	}

	drain_nats_error := backend_app.nats_connection.Drain()
	if drain_nats_error != nil {
		shutdown_error_list = append(shutdown_error_list, fmt.Errorf("failed_to_drain_nats_connection: %w", drain_nats_error))
	}
	backend_app.nats_connection.Close()
	backend_app.postgres_pool.Close()

	if len(shutdown_error_list) > 0 {
		return errors.Join(shutdown_error_list...)
	}

	return nil
}
