package storage

import (
	"context"
	"errors"
	"fmt"

	"ecochitas/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Truck_position_repository struct {
	postgres_pool *pgxpool.Pool
}

func New_truck_position_repository(postgres_pool *pgxpool.Pool) *Truck_position_repository {
	return &Truck_position_repository{
		postgres_pool: postgres_pool,
	}
}

func (truck_position_repository *Truck_position_repository) Save_truck_position(
	application_context context.Context,
	gps_location_event domain.Gps_location_event,
) error {
	insert_statement := `
		INSERT INTO truck_positions (
			truck_identifier,
			latitude,
			longitude,
			speed_kmh,
			heading_degrees,
			captured_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (truck_identifier, captured_at)
		DO UPDATE SET
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			speed_kmh = EXCLUDED.speed_kmh,
			heading_degrees = EXCLUDED.heading_degrees,
			received_at = NOW();
	`

	_, execute_insert_error := truck_position_repository.postgres_pool.Exec(
		application_context,
		insert_statement,
		gps_location_event.Truck_identifier,
		gps_location_event.Latitude,
		gps_location_event.Longitude,
		gps_location_event.Speed_kmh,
		gps_location_event.Heading_degrees,
		gps_location_event.Captured_at,
	)
	if execute_insert_error != nil {
		return fmt.Errorf("failed_to_insert_truck_position: %w", execute_insert_error)
	}

	return nil
}

func (truck_position_repository *Truck_position_repository) Get_latest_truck_position(
	application_context context.Context,
	truck_identifier string,
) (*domain.Truck_latest_position, error) {
	query_statement := `
		SELECT
			truck_identifier,
			latitude::float8,
			longitude::float8,
			speed_kmh::float8,
			heading_degrees::float8,
			captured_at,
			received_at
		FROM truck_positions
		WHERE truck_identifier = $1
		ORDER BY captured_at DESC
		LIMIT 1;
	`

	var latest_truck_position domain.Truck_latest_position
	query_error := truck_position_repository.postgres_pool.QueryRow(
		application_context,
		query_statement,
		truck_identifier,
	).Scan(
		&latest_truck_position.Truck_identifier,
		&latest_truck_position.Latitude,
		&latest_truck_position.Longitude,
		&latest_truck_position.Speed_kmh,
		&latest_truck_position.Heading_degrees,
		&latest_truck_position.Captured_at,
		&latest_truck_position.Received_at,
	)
	if query_error != nil {
		if errors.Is(query_error, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed_to_query_latest_truck_position: %w", query_error)
	}

	return &latest_truck_position, nil
}

func (truck_position_repository *Truck_position_repository) List_latest_truck_positions(
	application_context context.Context,
) ([]domain.Truck_latest_position, error) {
	query_statement := `
		SELECT DISTINCT ON (truck_identifier)
			truck_identifier,
			latitude::float8,
			longitude::float8,
			speed_kmh::float8,
			heading_degrees::float8,
			captured_at,
			received_at
		FROM truck_positions
		ORDER BY truck_identifier, captured_at DESC;
	`

	query_rows, query_error := truck_position_repository.postgres_pool.Query(application_context, query_statement)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_latest_truck_positions: %w", query_error)
	}
	defer query_rows.Close()

	latest_position_list := make([]domain.Truck_latest_position, 0)
	for query_rows.Next() {
		var latest_position_item domain.Truck_latest_position

		scan_error := query_rows.Scan(
			&latest_position_item.Truck_identifier,
			&latest_position_item.Latitude,
			&latest_position_item.Longitude,
			&latest_position_item.Speed_kmh,
			&latest_position_item.Heading_degrees,
			&latest_position_item.Captured_at,
			&latest_position_item.Received_at,
		)
		if scan_error != nil {
			return nil, fmt.Errorf("failed_to_scan_latest_truck_position_row: %w", scan_error)
		}

		latest_position_list = append(latest_position_list, latest_position_item)
	}

	rows_iteration_error := query_rows.Err()
	if rows_iteration_error != nil {
		return nil, fmt.Errorf("failed_during_latest_truck_positions_rows_iteration: %w", rows_iteration_error)
	}

	return latest_position_list, nil
}
