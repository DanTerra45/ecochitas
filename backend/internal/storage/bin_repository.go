package storage

import (
	"context"
	"fmt"

	"ecochitas/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Bin_repository struct {
	postgres_pool *pgxpool.Pool
}

func New_bin_repository(postgres_pool *pgxpool.Pool) *Bin_repository {
	return &Bin_repository{
		postgres_pool: postgres_pool,
	}
}

func (bin_repository *Bin_repository) List_bins(application_context context.Context) ([]domain.Bin_status, error) {
	query_statement := `
		SELECT
			id::text,
			bin_code,
			zone_name,
			latitude::float8,
			longitude::float8,
			fill_percentage,
			sensor_status,
			bin_status,
			last_emptied_at
		FROM bins
		ORDER BY zone_name, bin_code;
	`

	query_rows, query_error := bin_repository.postgres_pool.Query(application_context, query_statement)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_bins: %w", query_error)
	}
	defer query_rows.Close()

	bin_status_list := make([]domain.Bin_status, 0)
	for query_rows.Next() {
		var bin_status_item domain.Bin_status

		scan_error := query_rows.Scan(
			&bin_status_item.Bin_identifier,
			&bin_status_item.Bin_code,
			&bin_status_item.Zone_name,
			&bin_status_item.Latitude,
			&bin_status_item.Longitude,
			&bin_status_item.Fill_percentage,
			&bin_status_item.Sensor_status,
			&bin_status_item.Bin_status_label,
			&bin_status_item.Last_emptied_at,
		)
		if scan_error != nil {
			return nil, fmt.Errorf("failed_to_scan_bin_row: %w", scan_error)
		}

		bin_status_list = append(bin_status_list, bin_status_item)
	}

	rows_iteration_error := query_rows.Err()
	if rows_iteration_error != nil {
		return nil, fmt.Errorf("failed_during_bin_rows_iteration: %w", rows_iteration_error)
	}

	return bin_status_list, nil
}
