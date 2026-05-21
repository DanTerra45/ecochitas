package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"ecochitas/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	default_route_revision_list_limit        = 20
	max_route_revision_list_limit            = 200
	default_route_assignment_list_limit      = 50
	max_route_assignment_list_limit          = 200
	default_route_deviation_alert_list_limit = 50
	max_route_deviation_alert_list_limit     = 200
	default_deviation_threshold_meters       = 200.0
	max_deviation_threshold_meters           = 5000.0
	route_revision_change_type_created       = "route_created"
	route_revision_change_type_updated       = "route_updated"
	route_revision_change_type_stops_synced  = "route_stops_synced"
	routing_status_not_required              = "not_required"
	routing_status_routed                    = "routed"
	routing_status_fallback_straight_line    = "fallback_straight_line"
	routing_provider_osrm                    = "osrm"
)

type Route_repository struct {
	postgres_pool     *pgxpool.Pool
	route_path_router *osrm_route_service
}

type collection_route_row struct {
	Route_identifier          string
	Route_code                string
	Route_name                string
	Zone_name                 string
	Collection_weekday        int
	Is_active                 bool
	Stop_total                int
	Raw_road_path_coordinates []byte
	Routing_status            string
	Routing_provider          string
	Routed_at                 *time.Time
}

func New_route_repository(postgres_pool *pgxpool.Pool, osrm_base_url string) *Route_repository {
	return &Route_repository{
		postgres_pool:     postgres_pool,
		route_path_router: new_osrm_route_service(osrm_base_url),
	}
}

func (route_repository *Route_repository) List_collection_routes(
	application_context context.Context,
	collection_route_list_query domain.Collection_route_list_query,
) ([]domain.Collection_route_view, error) {
	if collection_route_list_query.Collection_weekday_filter < 0 || collection_route_list_query.Collection_weekday_filter > 7 {
		return nil, fmt.Errorf("invalid_collection_weekday_filter")
	}

	list_collection_routes_statement := `
		SELECT
			route.id::text,
			route.route_code,
			route.route_name,
			route.zone_name,
			route.collection_weekday::int,
			route.is_active,
			COALESCE(stop_total_subquery.stop_total, 0)::int,
			COALESCE(route.road_path_coordinates, '[]'::jsonb),
			COALESCE(route.routing_status, ''),
			COALESCE(route.routing_provider, ''),
			route.routed_at
		FROM collection_routes AS route
		LEFT JOIN LATERAL (
			SELECT COUNT(*)::int AS stop_total
			FROM route_stops AS stop
			WHERE stop.route_id = route.id
		) AS stop_total_subquery ON TRUE
		WHERE
			($1 = '' OR route.zone_name = $1) AND
			($2 = 0 OR route.collection_weekday = $2) AND
			($3 = FALSE OR route.is_active = $4)
		ORDER BY route.zone_name, route.route_code;
	`

	query_rows, query_error := route_repository.postgres_pool.Query(
		application_context,
		list_collection_routes_statement,
		strings.TrimSpace(collection_route_list_query.Zone_name_filter),
		collection_route_list_query.Collection_weekday_filter,
		collection_route_list_query.Has_is_active_filter,
		collection_route_list_query.Is_active_filter_value,
	)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_collection_routes: %w", query_error)
	}
	defer query_rows.Close()

	collection_route_view_list := make([]domain.Collection_route_view, 0)
	for query_rows.Next() {
		var route_row collection_route_row
		scan_row_error := query_rows.Scan(
			&route_row.Route_identifier,
			&route_row.Route_code,
			&route_row.Route_name,
			&route_row.Zone_name,
			&route_row.Collection_weekday,
			&route_row.Is_active,
			&route_row.Stop_total,
			&route_row.Raw_road_path_coordinates,
			&route_row.Routing_status,
			&route_row.Routing_provider,
			&route_row.Routed_at,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_collection_route_row: %w", scan_row_error)
		}

		route_stop_view_list, load_route_stops_error := load_route_stop_view_list_by_route_identifier(
			application_context,
			route_repository.postgres_pool,
			route_row.Route_identifier,
		)
		if load_route_stops_error != nil {
			return nil, load_route_stops_error
		}

		collection_route_view_item, map_collection_route_error := map_collection_route_row_to_view(
			route_row,
			route_stop_view_list,
		)
		if map_collection_route_error != nil {
			return nil, map_collection_route_error
		}
		collection_route_view_list = append(collection_route_view_list, collection_route_view_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_collection_routes_rows_iteration: %w", iterate_rows_error)
	}

	return collection_route_view_list, nil
}

func (route_repository *Route_repository) Create_collection_route(
	application_context context.Context,
	collection_route_create_command domain.Collection_route_create_command,
) (*domain.Collection_route_view, error) {
	normalized_route_code := strings.TrimSpace(collection_route_create_command.Route_code)
	normalized_route_name := strings.TrimSpace(collection_route_create_command.Route_name)
	normalized_zone_name := strings.TrimSpace(collection_route_create_command.Zone_name)

	if normalized_route_code == "" {
		return nil, fmt.Errorf("route_code_is_required")
	}
	if normalized_route_name == "" {
		return nil, fmt.Errorf("route_name_is_required")
	}
	if normalized_zone_name == "" {
		return nil, fmt.Errorf("zone_name_is_required")
	}
	if collection_route_create_command.Collection_weekday < 1 || collection_route_create_command.Collection_weekday > 7 {
		return nil, fmt.Errorf("collection_weekday_out_of_range")
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_create_collection_route_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	create_collection_route_statement := `
		INSERT INTO collection_routes (
			route_code,
			route_name,
			zone_name,
			collection_weekday,
			is_active
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING id::text;
	`

	var created_route_identifier string
	create_collection_route_error := transaction_handler.QueryRow(
		application_context,
		create_collection_route_statement,
		normalized_route_code,
		normalized_route_name,
		normalized_zone_name,
		collection_route_create_command.Collection_weekday,
		collection_route_create_command.Is_active,
	).Scan(&created_route_identifier)
	if create_collection_route_error != nil {
		if is_unique_constraint_error(create_collection_route_error, "collection_routes_route_code_key") {
			return nil, fmt.Errorf("route_code_already_exists")
		}
		return nil, fmt.Errorf("failed_to_create_collection_route: %w", create_collection_route_error)
	}

	route_change_payload := map[string]any{
		"route_code":         normalized_route_code,
		"route_name":         normalized_route_name,
		"zone_name":          normalized_zone_name,
		"collection_weekday": collection_route_create_command.Collection_weekday,
		"is_active":          collection_route_create_command.Is_active,
		"route_identifier":   created_route_identifier,
		"action_description": "collection_route_created",
	}
	insert_revision_error := insert_collection_route_revision(
		application_context,
		transaction_handler,
		created_route_identifier,
		route_revision_change_type_created,
		strings.TrimSpace(collection_route_create_command.Authenticated_user_identifier),
		route_change_payload,
	)
	if insert_revision_error != nil {
		return nil, insert_revision_error
	}

	refresh_route_road_path_error := route_repository.refresh_route_road_path_coordinates(
		application_context,
		transaction_handler,
		created_route_identifier,
		nil,
	)
	if refresh_route_road_path_error != nil {
		return nil, refresh_route_road_path_error
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_create_collection_route_transaction: %w", commit_transaction_error)
	}

	return load_collection_route_view_by_identifier(
		application_context,
		route_repository.postgres_pool,
		created_route_identifier,
	)
}

func (route_repository *Route_repository) Create_demo_route(
	application_context context.Context,
	demo_route_create_command domain.Demo_route_create_command,
) (*domain.Collection_route_view, error) {
	normalized_route_code := strings.TrimSpace(demo_route_create_command.Route_code)
	normalized_route_name := strings.TrimSpace(demo_route_create_command.Route_name)
	normalized_zone_name := strings.TrimSpace(demo_route_create_command.Zone_name)

	if normalized_route_code == "" {
		return nil, fmt.Errorf("route_code_is_required")
	}
	if normalized_route_name == "" {
		return nil, fmt.Errorf("route_name_is_required")
	}
	if normalized_zone_name == "" {
		return nil, fmt.Errorf("zone_name_is_required")
	}
	if demo_route_create_command.Collection_weekday < 1 || demo_route_create_command.Collection_weekday > 7 {
		return nil, fmt.Errorf("collection_weekday_out_of_range")
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_create_demo_route_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	create_collection_route_statement := `
		INSERT INTO collection_routes (
			route_code,
			route_name,
			zone_name,
			collection_weekday,
			is_active
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING id::text;
	`

	var created_route_identifier string
	create_collection_route_error := transaction_handler.QueryRow(
		application_context,
		create_collection_route_statement,
		normalized_route_code,
		normalized_route_name,
		normalized_zone_name,
		demo_route_create_command.Collection_weekday,
		demo_route_create_command.Is_active,
	).Scan(&created_route_identifier)
	if create_collection_route_error != nil {
		if is_unique_constraint_error(create_collection_route_error, "collection_routes_route_code_key") {
			return nil, fmt.Errorf("route_code_already_exists")
		}
		return nil, fmt.Errorf("failed_to_create_collection_route: %w", create_collection_route_error)
	}

	var raw_road_path_coordinates []byte
	routing_status := routing_status_not_required
	routing_provider := ""

	if len(demo_route_create_command.Points) >= 2 {
		route_path_coordinate_list := make([]domain.Route_path_coordinate, len(demo_route_create_command.Points))
		for i, point := range demo_route_create_command.Points {
			route_path_coordinate_list[i] = domain.Route_path_coordinate{
				Stop_order: i + 1,
				Latitude:   point.Latitude,
				Longitude:  point.Longitude,
			}
		}

		road_path_coordinate_list, osrm_error := route_repository.route_path_router.calculate_road_path_coordinates(
			application_context,
			route_path_coordinate_list,
		)

		if osrm_error == nil && len(road_path_coordinate_list) > 0 {
			marshaled_path, err := json.Marshal(road_path_coordinate_list)
			if err == nil {
				raw_road_path_coordinates = marshaled_path
				routing_status = routing_status_routed
				routing_provider = routing_provider_osrm
			}
		} else {
			// Fallback straight lines
			straight_line_path := build_straight_line_road_path_coordinates(route_path_coordinate_list)
			marshaled_path, _ := json.Marshal(straight_line_path)
			raw_road_path_coordinates = marshaled_path
			routing_status = routing_status_fallback_straight_line
		}
	}

	if len(raw_road_path_coordinates) > 0 {
		update_path_statement := `
			UPDATE collection_routes
			SET
				road_path_coordinates = $2::jsonb,
				routing_status = $3,
				routing_provider = $4,
				routed_at = NOW()
			WHERE id = $1::uuid;
		`
		_, update_path_error := transaction_handler.Exec(
			application_context,
			update_path_statement,
			created_route_identifier,
			raw_road_path_coordinates,
			routing_status,
			routing_provider,
		)
		if update_path_error != nil {
			return nil, fmt.Errorf("failed_to_update_demo_route_path: %w", update_path_error)
		}
	}

	route_change_payload := map[string]any{
		"route_code":         normalized_route_code,
		"route_name":         normalized_route_name,
		"zone_name":          normalized_zone_name,
		"collection_weekday": demo_route_create_command.Collection_weekday,
		"is_active":          demo_route_create_command.Is_active,
		"route_identifier":   created_route_identifier,
		"action_description": "demo_route_created",
	}
	insert_revision_error := insert_collection_route_revision(
		application_context,
		transaction_handler,
		created_route_identifier,
		route_revision_change_type_created,
		strings.TrimSpace(demo_route_create_command.Authenticated_user_identifier),
		route_change_payload,
	)
	if insert_revision_error != nil {
		return nil, insert_revision_error
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_create_demo_route_transaction: %w", commit_transaction_error)
	}

	return load_collection_route_view_by_identifier(
		application_context,
		route_repository.postgres_pool,
		created_route_identifier,
	)
}

func (route_repository *Route_repository) Update_collection_route(
	application_context context.Context,
	collection_route_update_command domain.Collection_route_update_command,
) (*domain.Collection_route_view, error) {
	normalized_route_identifier := strings.TrimSpace(collection_route_update_command.Route_identifier)
	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	has_any_field_update := collection_route_update_command.Route_code != nil ||
		collection_route_update_command.Route_name != nil ||
		collection_route_update_command.Zone_name != nil ||
		collection_route_update_command.Collection_weekday != nil ||
		collection_route_update_command.Is_active != nil
	if !has_any_field_update {
		return nil, fmt.Errorf("no_route_fields_to_update")
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_update_collection_route_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	load_existing_route_statement := `
		SELECT
			route_code,
			route_name,
			zone_name,
			collection_weekday::int,
			is_active
		FROM collection_routes
		WHERE id = $1::uuid
		FOR UPDATE;
	`

	var existing_route_code string
	var existing_route_name string
	var existing_zone_name string
	var existing_collection_weekday int
	var existing_is_active bool
	load_existing_route_error := transaction_handler.QueryRow(
		application_context,
		load_existing_route_statement,
		normalized_route_identifier,
	).Scan(
		&existing_route_code,
		&existing_route_name,
		&existing_zone_name,
		&existing_collection_weekday,
		&existing_is_active,
	)
	if load_existing_route_error != nil {
		if errors.Is(load_existing_route_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("route_not_found")
		}
		return nil, fmt.Errorf("failed_to_load_existing_collection_route: %w", load_existing_route_error)
	}

	updated_route_code := existing_route_code
	updated_route_name := existing_route_name
	updated_zone_name := existing_zone_name
	updated_collection_weekday := existing_collection_weekday
	updated_is_active := existing_is_active
	updated_field_name_list := make([]string, 0)

	if collection_route_update_command.Route_code != nil {
		normalized_route_code := strings.TrimSpace(*collection_route_update_command.Route_code)
		if normalized_route_code == "" {
			return nil, fmt.Errorf("route_code_is_required")
		}
		updated_route_code = normalized_route_code
		updated_field_name_list = append(updated_field_name_list, "route_code")
	}

	if collection_route_update_command.Route_name != nil {
		normalized_route_name := strings.TrimSpace(*collection_route_update_command.Route_name)
		if normalized_route_name == "" {
			return nil, fmt.Errorf("route_name_is_required")
		}
		updated_route_name = normalized_route_name
		updated_field_name_list = append(updated_field_name_list, "route_name")
	}

	if collection_route_update_command.Zone_name != nil {
		normalized_zone_name := strings.TrimSpace(*collection_route_update_command.Zone_name)
		if normalized_zone_name == "" {
			return nil, fmt.Errorf("zone_name_is_required")
		}
		updated_zone_name = normalized_zone_name
		updated_field_name_list = append(updated_field_name_list, "zone_name")
	}

	if collection_route_update_command.Collection_weekday != nil {
		if *collection_route_update_command.Collection_weekday < 1 || *collection_route_update_command.Collection_weekday > 7 {
			return nil, fmt.Errorf("collection_weekday_out_of_range")
		}
		updated_collection_weekday = *collection_route_update_command.Collection_weekday
		updated_field_name_list = append(updated_field_name_list, "collection_weekday")
	}

	if collection_route_update_command.Is_active != nil {
		updated_is_active = *collection_route_update_command.Is_active
		updated_field_name_list = append(updated_field_name_list, "is_active")
	}

	update_collection_route_statement := `
		UPDATE collection_routes
		SET
			route_code = $2,
			route_name = $3,
			zone_name = $4,
			collection_weekday = $5,
			is_active = $6,
			updated_at = NOW()
		WHERE id = $1::uuid;
	`
	_, update_collection_route_error := transaction_handler.Exec(
		application_context,
		update_collection_route_statement,
		normalized_route_identifier,
		updated_route_code,
		updated_route_name,
		updated_zone_name,
		updated_collection_weekday,
		updated_is_active,
	)
	if update_collection_route_error != nil {
		if is_unique_constraint_error(update_collection_route_error, "collection_routes_route_code_key") {
			return nil, fmt.Errorf("route_code_already_exists")
		}
		return nil, fmt.Errorf("failed_to_update_collection_route: %w", update_collection_route_error)
	}

	route_change_payload := map[string]any{
		"updated_fields": updated_field_name_list,
		"previous_values": map[string]any{
			"route_code":         existing_route_code,
			"route_name":         existing_route_name,
			"zone_name":          existing_zone_name,
			"collection_weekday": existing_collection_weekday,
			"is_active":          existing_is_active,
		},
		"current_values": map[string]any{
			"route_code":         updated_route_code,
			"route_name":         updated_route_name,
			"zone_name":          updated_zone_name,
			"collection_weekday": updated_collection_weekday,
			"is_active":          updated_is_active,
		},
	}
	insert_revision_error := insert_collection_route_revision(
		application_context,
		transaction_handler,
		normalized_route_identifier,
		route_revision_change_type_updated,
		strings.TrimSpace(collection_route_update_command.Authenticated_user_identifier),
		route_change_payload,
	)
	if insert_revision_error != nil {
		return nil, insert_revision_error
	}

	refresh_route_road_path_error := route_repository.refresh_route_road_path_coordinates(
		application_context,
		transaction_handler,
		normalized_route_identifier,
		nil,
	)
	if refresh_route_road_path_error != nil {
		return nil, refresh_route_road_path_error
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_update_collection_route_transaction: %w", commit_transaction_error)
	}

	return load_collection_route_view_by_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
}

func (route_repository *Route_repository) List_route_stops_by_route_identifier(
	application_context context.Context,
	route_identifier string,
) ([]domain.Route_stop_view, error) {
	normalized_route_identifier := strings.TrimSpace(route_identifier)
	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	route_exists, check_exists_error := route_exists_by_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
	if check_exists_error != nil {
		return nil, check_exists_error
	}
	if !route_exists {
		return nil, fmt.Errorf("route_not_found")
	}

	return load_route_stop_view_list_by_route_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
}

func (route_repository *Route_repository) Sync_route_stops(
	application_context context.Context,
	route_stop_sync_command domain.Route_stop_sync_command,
) ([]domain.Route_stop_view, error) {
	normalized_route_identifier := strings.TrimSpace(route_stop_sync_command.Route_identifier)
	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	normalized_route_stop_sync_item_list, normalize_items_error := normalize_route_stop_sync_item_list(
		route_stop_sync_command.Stop_item_list,
	)
	if normalize_items_error != nil {
		return nil, normalize_items_error
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_route_stop_sync_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	load_route_zone_statement := `
		SELECT zone_name
		FROM collection_routes
		WHERE id = $1::uuid
		FOR UPDATE;
	`
	var route_zone_name string
	load_route_zone_error := transaction_handler.QueryRow(
		application_context,
		load_route_zone_statement,
		normalized_route_identifier,
	).Scan(&route_zone_name)
	if load_route_zone_error != nil {
		if errors.Is(load_route_zone_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("route_not_found")
		}
		return nil, fmt.Errorf("failed_to_load_route_zone_name: %w", load_route_zone_error)
	}

	delete_existing_stops_statement := `
		DELETE FROM route_stops
		WHERE route_id = $1::uuid;
	`
	_, delete_existing_stops_error := transaction_handler.Exec(
		application_context,
		delete_existing_stops_statement,
		normalized_route_identifier,
	)
	if delete_existing_stops_error != nil {
		return nil, fmt.Errorf("failed_to_delete_existing_route_stops: %w", delete_existing_stops_error)
	}

	insert_route_stop_statement := `
		INSERT INTO route_stops (
			route_id,
			bin_id,
			stop_order,
			planned_time
		) VALUES (
			$1::uuid,
			$2::uuid,
			$3,
			NULLIF($4, '')::time
		);
	`
	load_bin_zone_statement := `
		SELECT zone_name
		FROM bins
		WHERE id = $1::uuid;
	`

	for _, route_stop_sync_item := range normalized_route_stop_sync_item_list {
		var bin_zone_name string
		load_bin_zone_error := transaction_handler.QueryRow(
			application_context,
			load_bin_zone_statement,
			route_stop_sync_item.Bin_identifier,
		).Scan(&bin_zone_name)
		if load_bin_zone_error != nil {
			if errors.Is(load_bin_zone_error, pgx.ErrNoRows) {
				return nil, fmt.Errorf("bin_not_found")
			}
			return nil, fmt.Errorf("failed_to_load_bin_zone_name: %w", load_bin_zone_error)
		}

		if strings.TrimSpace(bin_zone_name) != strings.TrimSpace(route_zone_name) {
			return nil, fmt.Errorf("bin_zone_mismatch")
		}

		_, insert_route_stop_error := transaction_handler.Exec(
			application_context,
			insert_route_stop_statement,
			normalized_route_identifier,
			route_stop_sync_item.Bin_identifier,
			route_stop_sync_item.Stop_order,
			route_stop_sync_item.Planned_time,
		)
		if insert_route_stop_error != nil {
			return nil, fmt.Errorf("failed_to_insert_route_stop: %w", insert_route_stop_error)
		}
	}

	route_stop_view_list, load_route_stop_view_list_error := load_route_stop_view_list_by_route_identifier(
		application_context,
		transaction_handler,
		normalized_route_identifier,
	)
	if load_route_stop_view_list_error != nil {
		return nil, load_route_stop_view_list_error
	}

	route_change_payload := map[string]any{
		"route_identifier": normalized_route_identifier,
		"stop_total":       len(route_stop_view_list),
		"stop_item_list":   normalized_route_stop_sync_item_list,
	}
	insert_revision_error := insert_collection_route_revision(
		application_context,
		transaction_handler,
		normalized_route_identifier,
		route_revision_change_type_stops_synced,
		strings.TrimSpace(route_stop_sync_command.Authenticated_user_identifier),
		route_change_payload,
	)
	if insert_revision_error != nil {
		return nil, insert_revision_error
	}

	refresh_route_road_path_error := route_repository.refresh_route_road_path_coordinates(
		application_context,
		transaction_handler,
		normalized_route_identifier,
		route_stop_view_list,
	)
	if refresh_route_road_path_error != nil {
		return nil, refresh_route_road_path_error
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_route_stop_sync_transaction: %w", commit_transaction_error)
	}

	return route_stop_view_list, nil
}

func (route_repository *Route_repository) List_collection_route_revisions(
	application_context context.Context,
	collection_route_revision_list_query domain.Collection_route_revision_list_query,
) (*domain.Collection_route_revision_list_result, error) {
	normalized_route_identifier := strings.TrimSpace(collection_route_revision_list_query.Route_identifier)
	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	route_exists, route_exists_error := route_exists_by_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
	if route_exists_error != nil {
		return nil, route_exists_error
	}
	if !route_exists {
		return nil, fmt.Errorf("route_not_found")
	}

	normalized_limit := collection_route_revision_list_query.Limit
	if normalized_limit <= 0 {
		normalized_limit = default_route_revision_list_limit
	}
	if normalized_limit > max_route_revision_list_limit {
		normalized_limit = max_route_revision_list_limit
	}

	list_collection_route_revisions_statement := `
		SELECT
			revision.id::text,
			revision.route_id::text,
			revision.revision_number,
			revision.change_type,
			COALESCE(revision.changed_by_user_identifier, ''),
			revision.change_payload,
			revision.created_at
		FROM collection_route_revisions AS revision
		WHERE revision.route_id = $1::uuid
		ORDER BY revision.revision_number DESC
		LIMIT $2;
	`

	query_rows, query_rows_error := route_repository.postgres_pool.Query(
		application_context,
		list_collection_route_revisions_statement,
		normalized_route_identifier,
		normalized_limit,
	)
	if query_rows_error != nil {
		return nil, fmt.Errorf("failed_to_query_collection_route_revisions: %w", query_rows_error)
	}
	defer query_rows.Close()

	revision_item_list := make([]domain.Collection_route_revision_record, 0)
	for query_rows.Next() {
		var revision_item domain.Collection_route_revision_record
		var raw_change_payload []byte
		scan_row_error := query_rows.Scan(
			&revision_item.Revision_identifier,
			&revision_item.Route_identifier,
			&revision_item.Revision_number,
			&revision_item.Change_type,
			&revision_item.Changed_by_user_identifier,
			&raw_change_payload,
			&revision_item.Created_at,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_collection_route_revision_row: %w", scan_row_error)
		}

		if len(raw_change_payload) == 0 {
			revision_item.Change_payload = map[string]any{}
		} else {
			unmarshal_change_payload_error := json.Unmarshal(raw_change_payload, &revision_item.Change_payload)
			if unmarshal_change_payload_error != nil {
				return nil, fmt.Errorf("failed_to_unmarshal_collection_route_revision_change_payload: %w", unmarshal_change_payload_error)
			}
		}

		revision_item.Changed_by_user_identifier = strings.TrimSpace(revision_item.Changed_by_user_identifier)
		revision_item_list = append(revision_item_list, revision_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_collection_route_revisions_rows_iteration: %w", iterate_rows_error)
	}

	return &domain.Collection_route_revision_list_result{
		Items: revision_item_list,
		Total: len(revision_item_list),
	}, nil
}

func (route_repository *Route_repository) Create_truck_route_assignment(
	application_context context.Context,
	truck_route_assignment_create_command domain.Truck_route_assignment_create_command,
) (*domain.Truck_route_assignment_view, error) {
	normalized_truck_identifier := strings.TrimSpace(truck_route_assignment_create_command.Truck_identifier)
	normalized_route_identifier := strings.TrimSpace(truck_route_assignment_create_command.Route_identifier)
	if normalized_truck_identifier == "" {
		return nil, fmt.Errorf("truck_identifier_is_required")
	}
	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	route_exists, route_exists_error := route_exists_by_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
	if route_exists_error != nil {
		return nil, route_exists_error
	}
	if !route_exists {
		return nil, fmt.Errorf("route_not_found")
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_create_truck_route_assignment_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	deactivate_existing_assignment_statement := `
		UPDATE truck_route_assignments
		SET
			is_active = FALSE,
			unassigned_at = NOW(),
			updated_at = NOW()
		WHERE truck_identifier = $1
			AND is_active = TRUE;
	`
	_, deactivate_existing_assignment_error := transaction_handler.Exec(
		application_context,
		deactivate_existing_assignment_statement,
		normalized_truck_identifier,
	)
	if deactivate_existing_assignment_error != nil {
		return nil, fmt.Errorf("failed_to_deactivate_existing_truck_route_assignment: %w", deactivate_existing_assignment_error)
	}

	create_assignment_statement := `
		INSERT INTO truck_route_assignments (
			truck_identifier,
			route_id,
			is_active,
			assigned_by_user_identifier,
			assignment_notes,
			assigned_at
		) VALUES (
			$1,
			$2::uuid,
			TRUE,
			NULLIF($3, ''),
			NULLIF($4, ''),
			NOW()
		)
		RETURNING id::text;
	`

	var created_assignment_identifier string
	create_assignment_error := transaction_handler.QueryRow(
		application_context,
		create_assignment_statement,
		normalized_truck_identifier,
		normalized_route_identifier,
		strings.TrimSpace(truck_route_assignment_create_command.Assigned_by_user_identifier),
		strings.TrimSpace(truck_route_assignment_create_command.Assignment_notes),
	).Scan(&created_assignment_identifier)
	if create_assignment_error != nil {
		return nil, fmt.Errorf("failed_to_create_truck_route_assignment: %w", create_assignment_error)
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_create_truck_route_assignment_transaction: %w", commit_transaction_error)
	}

	return load_truck_route_assignment_by_identifier(
		application_context,
		route_repository.postgres_pool,
		created_assignment_identifier,
	)
}

func (route_repository *Route_repository) List_truck_route_assignments(
	application_context context.Context,
	truck_route_assignment_list_query domain.Truck_route_assignment_list_query,
) (*domain.Truck_route_assignment_list_result, error) {
	normalized_route_identifier := strings.TrimSpace(truck_route_assignment_list_query.Route_identifier)
	if normalized_route_identifier != "" && !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	normalized_limit := truck_route_assignment_list_query.Limit
	if normalized_limit <= 0 {
		normalized_limit = default_route_assignment_list_limit
	}
	if normalized_limit > max_route_assignment_list_limit {
		normalized_limit = max_route_assignment_list_limit
	}

	has_is_active_filter := false
	is_active_filter_value := false
	if truck_route_assignment_list_query.Is_active_filter != nil {
		has_is_active_filter = true
		is_active_filter_value = *truck_route_assignment_list_query.Is_active_filter
	}

	list_assignments_statement := `
		SELECT
			assignment.id::text,
			assignment.truck_identifier,
			route.id::text,
			route.route_code,
			route.route_name,
			route.zone_name,
			route.collection_weekday::int,
			assignment.is_active,
			COALESCE(assignment.assigned_by_user_identifier, ''),
			COALESCE(assignment.assignment_notes, ''),
			assignment.assigned_at,
			assignment.unassigned_at
		FROM truck_route_assignments AS assignment
		INNER JOIN collection_routes AS route
			ON route.id = assignment.route_id
		WHERE
			($1 = '' OR assignment.truck_identifier = $1) AND
			($2 = '' OR assignment.route_id = $2::uuid) AND
			($3 = FALSE OR assignment.is_active = $4)
		ORDER BY assignment.assigned_at DESC
		LIMIT $5;
	`

	query_rows, query_rows_error := route_repository.postgres_pool.Query(
		application_context,
		list_assignments_statement,
		strings.TrimSpace(truck_route_assignment_list_query.Truck_identifier),
		normalized_route_identifier,
		has_is_active_filter,
		is_active_filter_value,
		normalized_limit,
	)
	if query_rows_error != nil {
		return nil, fmt.Errorf("failed_to_query_truck_route_assignments: %w", query_rows_error)
	}
	defer query_rows.Close()

	assignment_item_list := make([]domain.Truck_route_assignment_view, 0)
	for query_rows.Next() {
		var assignment_item domain.Truck_route_assignment_view
		scan_row_error := query_rows.Scan(
			&assignment_item.Assignment_identifier,
			&assignment_item.Truck_identifier,
			&assignment_item.Route_identifier,
			&assignment_item.Route_code,
			&assignment_item.Route_name,
			&assignment_item.Zone_name,
			&assignment_item.Collection_weekday,
			&assignment_item.Is_active,
			&assignment_item.Assigned_by_user_identifier,
			&assignment_item.Assignment_notes,
			&assignment_item.Assigned_at,
			&assignment_item.Unassigned_at,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_truck_route_assignment_row: %w", scan_row_error)
		}

		assignment_item.Assigned_by_user_identifier = strings.TrimSpace(assignment_item.Assigned_by_user_identifier)
		assignment_item.Assignment_notes = strings.TrimSpace(assignment_item.Assignment_notes)
		assignment_item_list = append(assignment_item_list, assignment_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_truck_route_assignments_rows_iteration: %w", iterate_rows_error)
	}

	return &domain.Truck_route_assignment_list_result{
		Items: assignment_item_list,
		Total: len(assignment_item_list),
	}, nil
}

func (route_repository *Route_repository) Create_route_deviation_alert(
	application_context context.Context,
	route_deviation_alert_create_command domain.Route_deviation_alert_create_command,
) (*domain.Route_deviation_alert_record, error) {
	truck_route_deviation_view, get_truck_route_deviation_error := route_repository.Get_truck_route_deviation(
		application_context,
		domain.Truck_route_deviation_query{
			Truck_identifier:           route_deviation_alert_create_command.Truck_identifier,
			Route_identifier:           route_deviation_alert_create_command.Route_identifier,
			Deviation_threshold_meters: route_deviation_alert_create_command.Deviation_threshold_meters,
		},
	)
	if get_truck_route_deviation_error != nil {
		return nil, get_truck_route_deviation_error
	}
	if !truck_route_deviation_view.Is_off_route {
		return nil, fmt.Errorf("truck_is_within_route_threshold")
	}

	derived_severity_level := derive_route_deviation_severity_level(
		truck_route_deviation_view.Distance_to_route_meters,
		truck_route_deviation_view.Deviation_threshold_meters,
	)
	metadata_payload := map[string]any{
		"captured_at":                   truck_route_deviation_view.Captured_at,
		"truck_latitude":                truck_route_deviation_view.Truck_latitude,
		"truck_longitude":               truck_route_deviation_view.Truck_longitude,
		"nearest_route_stop_identifier": truck_route_deviation_view.Nearest_route_stop_identifier,
		"nearest_stop_order":            truck_route_deviation_view.Nearest_stop_order,
		"nearest_bin_identifier":        truck_route_deviation_view.Nearest_bin_identifier,
		"nearest_bin_code":              truck_route_deviation_view.Nearest_bin_code,
		"nearest_bin_latitude":          truck_route_deviation_view.Nearest_bin_latitude,
		"nearest_bin_longitude":         truck_route_deviation_view.Nearest_bin_longitude,
	}
	serialized_metadata_payload, marshal_metadata_payload_error := json.Marshal(metadata_payload)
	if marshal_metadata_payload_error != nil {
		return nil, fmt.Errorf("failed_to_marshal_route_deviation_alert_metadata_payload: %w", marshal_metadata_payload_error)
	}

	transaction_handler, begin_transaction_error := route_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_create_route_deviation_alert_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	create_alert_statement := `
		INSERT INTO route_deviation_alert_events (
			truck_identifier,
			route_id,
			route_stop_id,
			distance_to_route_meters,
			deviation_threshold_meters,
			severity_level,
			alert_status,
			alert_notes,
			triggered_by_user_identifier,
			detected_at,
			metadata_payload
		) VALUES (
			$1,
			$2::uuid,
			$3::uuid,
			$4,
			$5,
			$6,
			'open',
			NULLIF($7, ''),
			NULLIF($8, ''),
			NOW(),
			$9::jsonb
		)
		RETURNING
			id::text,
			truck_identifier,
			route_id::text,
			COALESCE(route_stop_id::text, ''),
			distance_to_route_meters::float8,
			deviation_threshold_meters::float8,
			severity_level,
			alert_status,
			COALESCE(alert_notes, ''),
			COALESCE(triggered_by_user_identifier, ''),
			detected_at,
			resolved_at,
			metadata_payload,
			created_at;
	`

	created_route_deviation_alert_record := &domain.Route_deviation_alert_record{}
	var raw_metadata_payload []byte
	create_alert_error := transaction_handler.QueryRow(
		application_context,
		create_alert_statement,
		strings.TrimSpace(truck_route_deviation_view.Truck_identifier),
		strings.TrimSpace(truck_route_deviation_view.Route_identifier),
		strings.TrimSpace(truck_route_deviation_view.Nearest_route_stop_identifier),
		truck_route_deviation_view.Distance_to_route_meters,
		truck_route_deviation_view.Deviation_threshold_meters,
		derived_severity_level,
		strings.TrimSpace(route_deviation_alert_create_command.Alert_notes),
		strings.TrimSpace(route_deviation_alert_create_command.Triggered_by_user_identifier),
		serialized_metadata_payload,
	).Scan(
		&created_route_deviation_alert_record.Alert_identifier,
		&created_route_deviation_alert_record.Truck_identifier,
		&created_route_deviation_alert_record.Route_identifier,
		&created_route_deviation_alert_record.Route_stop_identifier,
		&created_route_deviation_alert_record.Distance_to_route_meters,
		&created_route_deviation_alert_record.Deviation_threshold_meters,
		&created_route_deviation_alert_record.Severity_level,
		&created_route_deviation_alert_record.Alert_status,
		&created_route_deviation_alert_record.Alert_notes,
		&created_route_deviation_alert_record.Triggered_by_user_identifier,
		&created_route_deviation_alert_record.Detected_at,
		&created_route_deviation_alert_record.Resolved_at,
		&raw_metadata_payload,
		&created_route_deviation_alert_record.Created_at,
	)
	if create_alert_error != nil {
		return nil, fmt.Errorf("failed_to_create_route_deviation_alert: %w", create_alert_error)
	}

	if len(raw_metadata_payload) == 0 {
		created_route_deviation_alert_record.Metadata_payload = map[string]any{}
	} else {
		unmarshal_metadata_payload_error := json.Unmarshal(
			raw_metadata_payload,
			&created_route_deviation_alert_record.Metadata_payload,
		)
		if unmarshal_metadata_payload_error != nil {
			return nil, fmt.Errorf("failed_to_unmarshal_route_deviation_alert_metadata_payload: %w", unmarshal_metadata_payload_error)
		}
	}

	insert_notification_event_statement := `
		INSERT INTO notification_outbox_events (
			event_type,
			channel_type,
			target_reference,
			title_text,
			body_text,
			payload
		) VALUES (
			'route_deviation_alert_recorded',
			'push',
			$1,
			$2,
			$3,
			$4::jsonb
		);
	`
	_, insert_notification_event_error := transaction_handler.Exec(
		application_context,
		insert_notification_event_statement,
		strings.TrimSpace(truck_route_deviation_view.Route_identifier),
		"Desvio detectado en ruta",
		"Camion "+strings.TrimSpace(truck_route_deviation_view.Truck_identifier)+" fuera de umbral",
		serialized_metadata_payload,
	)
	if insert_notification_event_error != nil {
		return nil, fmt.Errorf("failed_to_insert_route_deviation_notification_outbox_event: %w", insert_notification_event_error)
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_create_route_deviation_alert_transaction: %w", commit_transaction_error)
	}

	created_route_deviation_alert_record.Route_stop_identifier = strings.TrimSpace(
		created_route_deviation_alert_record.Route_stop_identifier,
	)
	created_route_deviation_alert_record.Alert_notes = strings.TrimSpace(created_route_deviation_alert_record.Alert_notes)
	created_route_deviation_alert_record.Triggered_by_user_identifier = strings.TrimSpace(
		created_route_deviation_alert_record.Triggered_by_user_identifier,
	)

	return created_route_deviation_alert_record, nil
}

func (route_repository *Route_repository) List_route_deviation_alerts(
	application_context context.Context,
	route_deviation_alert_list_query domain.Route_deviation_alert_list_query,
) (*domain.Route_deviation_alert_list_result, error) {
	normalized_route_identifier := strings.TrimSpace(route_deviation_alert_list_query.Route_identifier)
	if normalized_route_identifier != "" && !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	normalized_alert_status, normalize_alert_status_error := normalize_route_deviation_alert_status_filter(
		route_deviation_alert_list_query.Alert_status,
	)
	if normalize_alert_status_error != nil {
		return nil, normalize_alert_status_error
	}

	normalized_limit := route_deviation_alert_list_query.Limit
	if normalized_limit <= 0 {
		normalized_limit = default_route_deviation_alert_list_limit
	}
	if normalized_limit > max_route_deviation_alert_list_limit {
		normalized_limit = max_route_deviation_alert_list_limit
	}

	list_route_deviation_alerts_statement := `
		SELECT
			alert.id::text,
			alert.truck_identifier,
			alert.route_id::text,
			COALESCE(alert.route_stop_id::text, ''),
			alert.distance_to_route_meters::float8,
			alert.deviation_threshold_meters::float8,
			alert.severity_level,
			alert.alert_status,
			COALESCE(alert.alert_notes, ''),
			COALESCE(alert.triggered_by_user_identifier, ''),
			alert.detected_at,
			alert.resolved_at,
			alert.metadata_payload,
			alert.created_at
		FROM route_deviation_alert_events AS alert
		WHERE
			($1 = '' OR alert.truck_identifier = $1) AND
			($2 = '' OR alert.route_id = $2::uuid) AND
			($3 = '' OR alert.alert_status = $3)
		ORDER BY alert.detected_at DESC, alert.created_at DESC
		LIMIT $4;
	`

	query_rows, query_rows_error := route_repository.postgres_pool.Query(
		application_context,
		list_route_deviation_alerts_statement,
		strings.TrimSpace(route_deviation_alert_list_query.Truck_identifier),
		normalized_route_identifier,
		normalized_alert_status,
		normalized_limit,
	)
	if query_rows_error != nil {
		return nil, fmt.Errorf("failed_to_query_route_deviation_alerts: %w", query_rows_error)
	}
	defer query_rows.Close()

	route_deviation_alert_item_list := make([]domain.Route_deviation_alert_record, 0)
	for query_rows.Next() {
		var route_deviation_alert_item domain.Route_deviation_alert_record
		var raw_metadata_payload []byte
		scan_row_error := query_rows.Scan(
			&route_deviation_alert_item.Alert_identifier,
			&route_deviation_alert_item.Truck_identifier,
			&route_deviation_alert_item.Route_identifier,
			&route_deviation_alert_item.Route_stop_identifier,
			&route_deviation_alert_item.Distance_to_route_meters,
			&route_deviation_alert_item.Deviation_threshold_meters,
			&route_deviation_alert_item.Severity_level,
			&route_deviation_alert_item.Alert_status,
			&route_deviation_alert_item.Alert_notes,
			&route_deviation_alert_item.Triggered_by_user_identifier,
			&route_deviation_alert_item.Detected_at,
			&route_deviation_alert_item.Resolved_at,
			&raw_metadata_payload,
			&route_deviation_alert_item.Created_at,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_route_deviation_alert_row: %w", scan_row_error)
		}

		if len(raw_metadata_payload) == 0 {
			route_deviation_alert_item.Metadata_payload = map[string]any{}
		} else {
			unmarshal_metadata_payload_error := json.Unmarshal(
				raw_metadata_payload,
				&route_deviation_alert_item.Metadata_payload,
			)
			if unmarshal_metadata_payload_error != nil {
				return nil, fmt.Errorf("failed_to_unmarshal_route_deviation_alert_metadata_payload: %w", unmarshal_metadata_payload_error)
			}
		}

		route_deviation_alert_item.Route_stop_identifier = strings.TrimSpace(route_deviation_alert_item.Route_stop_identifier)
		route_deviation_alert_item.Alert_notes = strings.TrimSpace(route_deviation_alert_item.Alert_notes)
		route_deviation_alert_item.Triggered_by_user_identifier = strings.TrimSpace(
			route_deviation_alert_item.Triggered_by_user_identifier,
		)

		route_deviation_alert_item_list = append(route_deviation_alert_item_list, route_deviation_alert_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_route_deviation_alerts_rows_iteration: %w", iterate_rows_error)
	}

	return &domain.Route_deviation_alert_list_result{
		Items: route_deviation_alert_item_list,
		Total: len(route_deviation_alert_item_list),
	}, nil
}

func (route_repository *Route_repository) Get_truck_route_deviation(
	application_context context.Context,
	truck_route_deviation_query domain.Truck_route_deviation_query,
) (*domain.Truck_route_deviation_view, error) {
	normalized_truck_identifier := strings.TrimSpace(truck_route_deviation_query.Truck_identifier)
	normalized_route_identifier := strings.TrimSpace(truck_route_deviation_query.Route_identifier)
	if normalized_truck_identifier == "" {
		return nil, fmt.Errorf("truck_identifier_is_required")
	}

	if normalized_route_identifier == "" {
		resolved_route_identifier, resolve_route_identifier_error := resolve_active_route_identifier_by_truck_identifier(
			application_context,
			route_repository.postgres_pool,
			normalized_truck_identifier,
		)
		if resolve_route_identifier_error != nil {
			return nil, resolve_route_identifier_error
		}
		normalized_route_identifier = resolved_route_identifier
	}

	if !is_valid_uuid_string(normalized_route_identifier) {
		return nil, fmt.Errorf("invalid_route_identifier")
	}

	route_exists, route_exists_error := route_exists_by_identifier(
		application_context,
		route_repository.postgres_pool,
		normalized_route_identifier,
	)
	if route_exists_error != nil {
		return nil, route_exists_error
	}
	if !route_exists {
		return nil, fmt.Errorf("route_not_found")
	}

	deviation_threshold_meters := truck_route_deviation_query.Deviation_threshold_meters
	if deviation_threshold_meters <= 0 {
		deviation_threshold_meters = default_deviation_threshold_meters
	}
	if deviation_threshold_meters > max_deviation_threshold_meters {
		deviation_threshold_meters = max_deviation_threshold_meters
	}

	load_latest_truck_position_statement := `
		SELECT
			latitude::float8,
			longitude::float8,
			captured_at
		FROM truck_positions
		WHERE truck_identifier = $1
		ORDER BY captured_at DESC
		LIMIT 1;
	`

	var truck_latitude float64
	var truck_longitude float64
	var truck_position_captured_at time.Time
	load_latest_truck_position_error := route_repository.postgres_pool.QueryRow(
		application_context,
		load_latest_truck_position_statement,
		normalized_truck_identifier,
	).Scan(
		&truck_latitude,
		&truck_longitude,
		&truck_position_captured_at,
	)
	if load_latest_truck_position_error != nil {
		if errors.Is(load_latest_truck_position_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("truck_position_not_found")
		}
		return nil, fmt.Errorf("failed_to_load_latest_truck_position_for_route_deviation: %w", load_latest_truck_position_error)
	}

	load_nearest_route_stop_statement := `
		SELECT
			stop.id::text,
			stop.stop_order,
			bin.id::text,
			bin.bin_code,
			bin.latitude::float8,
			bin.longitude::float8,
			(
				6371000 * 2 * ASIN(
					SQRT(
						POWER(SIN(RADIANS(($1 - bin.latitude::float8) / 2)), 2) +
						COS(RADIANS($1)) * COS(RADIANS(bin.latitude::float8)) *
						POWER(SIN(RADIANS(($2 - bin.longitude::float8) / 2)), 2)
					)
				)
			) AS distance_meters
		FROM route_stops AS stop
		INNER JOIN bins AS bin
			ON bin.id = stop.bin_id
		WHERE stop.route_id = $3::uuid
		ORDER BY distance_meters ASC, stop.stop_order ASC
		LIMIT 1;
	`

	truck_route_deviation_view := &domain.Truck_route_deviation_view{
		Truck_identifier:           normalized_truck_identifier,
		Route_identifier:           normalized_route_identifier,
		Captured_at:                truck_position_captured_at.UTC(),
		Truck_latitude:             truck_latitude,
		Truck_longitude:            truck_longitude,
		Deviation_threshold_meters: deviation_threshold_meters,
	}
	load_nearest_route_stop_error := route_repository.postgres_pool.QueryRow(
		application_context,
		load_nearest_route_stop_statement,
		truck_latitude,
		truck_longitude,
		normalized_route_identifier,
	).Scan(
		&truck_route_deviation_view.Nearest_route_stop_identifier,
		&truck_route_deviation_view.Nearest_stop_order,
		&truck_route_deviation_view.Nearest_bin_identifier,
		&truck_route_deviation_view.Nearest_bin_code,
		&truck_route_deviation_view.Nearest_bin_latitude,
		&truck_route_deviation_view.Nearest_bin_longitude,
		&truck_route_deviation_view.Distance_to_route_meters,
	)
	if load_nearest_route_stop_error != nil {
		if errors.Is(load_nearest_route_stop_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("route_has_no_stops")
		}
		return nil, fmt.Errorf("failed_to_load_nearest_route_stop_for_route_deviation: %w", load_nearest_route_stop_error)
	}

	truck_route_deviation_view.Distance_to_route_meters = round_float_value(
		truck_route_deviation_view.Distance_to_route_meters,
		2,
	)
	truck_route_deviation_view.Is_off_route = truck_route_deviation_view.Distance_to_route_meters >
		deviation_threshold_meters

	return truck_route_deviation_view, nil
}

func route_exists_by_identifier(
	application_context context.Context,
	query_executor postgres_query_executor,
	route_identifier string,
) (bool, error) {
	query_route_exists_statement := `
		SELECT EXISTS (
			SELECT 1
			FROM collection_routes
			WHERE id = $1::uuid
		);
	`

	var route_exists bool
	query_route_exists_error := query_executor.QueryRow(
		application_context,
		query_route_exists_statement,
		route_identifier,
	).Scan(&route_exists)
	if query_route_exists_error != nil {
		return false, fmt.Errorf("failed_to_query_route_exists_by_identifier: %w", query_route_exists_error)
	}

	return route_exists, nil
}

func load_collection_route_view_by_identifier(
	application_context context.Context,
	query_executor postgres_query_executor,
	route_identifier string,
) (*domain.Collection_route_view, error) {
	load_collection_route_statement := `
		SELECT
			route.id::text,
			route.route_code,
			route.route_name,
			route.zone_name,
			route.collection_weekday::int,
			route.is_active,
			COALESCE(stop_total_subquery.stop_total, 0)::int,
			COALESCE(route.road_path_coordinates, '[]'::jsonb),
			COALESCE(route.routing_status, ''),
			COALESCE(route.routing_provider, ''),
			route.routed_at
		FROM collection_routes AS route
		LEFT JOIN LATERAL (
			SELECT COUNT(*)::int AS stop_total
			FROM route_stops AS stop
			WHERE stop.route_id = route.id
		) AS stop_total_subquery ON TRUE
		WHERE route.id = $1::uuid;
	`

	var route_row collection_route_row
	load_collection_route_error := query_executor.QueryRow(
		application_context,
		load_collection_route_statement,
		route_identifier,
	).Scan(
		&route_row.Route_identifier,
		&route_row.Route_code,
		&route_row.Route_name,
		&route_row.Zone_name,
		&route_row.Collection_weekday,
		&route_row.Is_active,
		&route_row.Stop_total,
		&route_row.Raw_road_path_coordinates,
		&route_row.Routing_status,
		&route_row.Routing_provider,
		&route_row.Routed_at,
	)
	if load_collection_route_error != nil {
		if errors.Is(load_collection_route_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("route_not_found")
		}
		return nil, fmt.Errorf("failed_to_load_collection_route_by_identifier: %w", load_collection_route_error)
	}

	route_stop_view_list, load_route_stops_error := load_route_stop_view_list_by_route_identifier(
		application_context,
		query_executor,
		route_row.Route_identifier,
	)
	if load_route_stops_error != nil {
		return nil, load_route_stops_error
	}

	collection_route_view, map_collection_route_error := map_collection_route_row_to_view(
		route_row,
		route_stop_view_list,
	)
	if map_collection_route_error != nil {
		return nil, map_collection_route_error
	}
	return &collection_route_view, nil
}

func load_route_stop_view_list_by_route_identifier(
	application_context context.Context,
	query_executor postgres_query_executor,
	route_identifier string,
) ([]domain.Route_stop_view, error) {
	list_route_stops_statement := `
		SELECT
			stop.id::text,
			stop.route_id::text,
			bin.id::text,
			bin.bin_code,
			bin.zone_name,
			stop.stop_order,
			COALESCE(stop.planned_time::text, ''),
			bin.latitude::float8,
			bin.longitude::float8
		FROM route_stops AS stop
		INNER JOIN bins AS bin
			ON bin.id = stop.bin_id
		WHERE stop.route_id = $1::uuid
		ORDER BY stop.stop_order;
	`

	query_rows, query_error := query_executor.Query(
		application_context,
		list_route_stops_statement,
		route_identifier,
	)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_route_stops: %w", query_error)
	}
	defer query_rows.Close()

	route_stop_view_list := make([]domain.Route_stop_view, 0)
	for query_rows.Next() {
		var route_stop_view_item domain.Route_stop_view
		scan_row_error := query_rows.Scan(
			&route_stop_view_item.Route_stop_identifier,
			&route_stop_view_item.Route_identifier,
			&route_stop_view_item.Bin_identifier,
			&route_stop_view_item.Bin_code,
			&route_stop_view_item.Zone_name,
			&route_stop_view_item.Stop_order,
			&route_stop_view_item.Planned_time,
			&route_stop_view_item.Latitude,
			&route_stop_view_item.Longitude,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_route_stop_row: %w", scan_row_error)
		}

		route_stop_view_item.Planned_time = strings.TrimSpace(route_stop_view_item.Planned_time)
		route_stop_view_list = append(route_stop_view_list, route_stop_view_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_route_stops_rows_iteration: %w", iterate_rows_error)
	}

	return route_stop_view_list, nil
}

func insert_collection_route_revision(
	application_context context.Context,
	query_executor postgres_query_executor,
	route_identifier string,
	change_type string,
	changed_by_user_identifier string,
	change_payload map[string]any,
) error {
	if change_payload == nil {
		change_payload = map[string]any{}
	}

	serialized_change_payload, marshal_change_payload_error := json.Marshal(change_payload)
	if marshal_change_payload_error != nil {
		return fmt.Errorf("failed_to_marshal_collection_route_revision_change_payload: %w", marshal_change_payload_error)
	}

	insert_collection_route_revision_statement := `
		WITH next_revision AS (
			SELECT
				COALESCE(MAX(revision_number), 0) + 1 AS revision_number
			FROM collection_route_revisions
			WHERE route_id = $1::uuid
		)
		INSERT INTO collection_route_revisions (
			route_id,
			revision_number,
			change_type,
			changed_by_user_identifier,
			change_payload
		)
		SELECT
			$1::uuid,
			next_revision.revision_number,
			$2,
			NULLIF($3, ''),
			$4::jsonb
		FROM next_revision;
	`

	_, insert_collection_route_revision_error := query_executor.Exec(
		application_context,
		insert_collection_route_revision_statement,
		route_identifier,
		change_type,
		strings.TrimSpace(changed_by_user_identifier),
		serialized_change_payload,
	)
	if insert_collection_route_revision_error != nil {
		return fmt.Errorf("failed_to_insert_collection_route_revision: %w", insert_collection_route_revision_error)
	}

	return nil
}

func normalize_route_stop_sync_item_list(
	raw_route_stop_sync_item_list []domain.Route_stop_sync_item,
) ([]domain.Route_stop_sync_item, error) {
	stop_order_seen_set := make(map[int]struct{}, len(raw_route_stop_sync_item_list))
	bin_identifier_seen_set := make(map[string]struct{}, len(raw_route_stop_sync_item_list))

	normalized_route_stop_sync_item_list := make([]domain.Route_stop_sync_item, 0, len(raw_route_stop_sync_item_list))
	for _, raw_route_stop_sync_item := range raw_route_stop_sync_item_list {
		normalized_bin_identifier := strings.TrimSpace(raw_route_stop_sync_item.Bin_identifier)
		if !is_valid_uuid_string(normalized_bin_identifier) {
			return nil, fmt.Errorf("invalid_bin_identifier")
		}
		if raw_route_stop_sync_item.Stop_order <= 0 {
			return nil, fmt.Errorf("route_stop_order_must_be_positive")
		}
		if _, already_exists := stop_order_seen_set[raw_route_stop_sync_item.Stop_order]; already_exists {
			return nil, fmt.Errorf("duplicate_route_stop_order")
		}
		if _, already_exists := bin_identifier_seen_set[normalized_bin_identifier]; already_exists {
			return nil, fmt.Errorf("duplicate_bin_identifier_in_route_stops")
		}

		normalized_planned_time, normalize_planned_time_error := normalize_route_stop_planned_time(
			raw_route_stop_sync_item.Planned_time,
		)
		if normalize_planned_time_error != nil {
			return nil, normalize_planned_time_error
		}

		stop_order_seen_set[raw_route_stop_sync_item.Stop_order] = struct{}{}
		bin_identifier_seen_set[normalized_bin_identifier] = struct{}{}
		normalized_route_stop_sync_item_list = append(
			normalized_route_stop_sync_item_list,
			domain.Route_stop_sync_item{
				Bin_identifier: normalized_bin_identifier,
				Stop_order:     raw_route_stop_sync_item.Stop_order,
				Planned_time:   normalized_planned_time,
			},
		)
	}

	return normalized_route_stop_sync_item_list, nil
}

func normalize_route_stop_planned_time(raw_planned_time string) (string, error) {
	normalized_planned_time := strings.TrimSpace(raw_planned_time)
	if normalized_planned_time == "" {
		return "", nil
	}

	parsed_planned_time_without_seconds, parse_without_seconds_error := time.Parse("15:04", normalized_planned_time)
	if parse_without_seconds_error == nil {
		return parsed_planned_time_without_seconds.Format("15:04:05"), nil
	}

	parsed_planned_time_with_seconds, parse_with_seconds_error := time.Parse("15:04:05", normalized_planned_time)
	if parse_with_seconds_error == nil {
		return parsed_planned_time_with_seconds.Format("15:04:05"), nil
	}

	return "", fmt.Errorf("invalid_planned_time")
}

func map_collection_route_row_to_view(
	route_row collection_route_row,
	route_stop_view_list []domain.Route_stop_view,
) (domain.Collection_route_view, error) {
	road_path_coordinate_list := make([]domain.Route_road_path_coordinate, 0)
	if len(route_row.Raw_road_path_coordinates) > 0 {
		unmarshal_road_path_error := json.Unmarshal(
			route_row.Raw_road_path_coordinates,
			&road_path_coordinate_list,
		)
		if unmarshal_road_path_error != nil {
			return domain.Collection_route_view{}, fmt.Errorf(
				"failed_to_unmarshal_route_road_path_coordinates: %w",
				unmarshal_road_path_error,
			)
		}
	}

	collection_route_view := domain.Collection_route_view{
		Route_identifier:      route_row.Route_identifier,
		Route_code:            route_row.Route_code,
		Route_name:            route_row.Route_name,
		Zone_name:             route_row.Zone_name,
		Collection_weekday:    route_row.Collection_weekday,
		Is_active:             route_row.Is_active,
		Stop_total:            route_row.Stop_total,
		Path_coordinates:      make([]domain.Route_path_coordinate, 0, len(route_stop_view_list)),
		Road_path_coordinates: road_path_coordinate_list,
		Routing_status:        strings.TrimSpace(route_row.Routing_status),
		Routing_provider:      strings.TrimSpace(route_row.Routing_provider),
	}
	if route_row.Routed_at != nil {
		routed_at_utc := route_row.Routed_at.UTC()
		collection_route_view.Routed_at = &routed_at_utc
	}

	for _, route_stop_view_item := range route_stop_view_list {
		collection_route_view.Path_coordinates = append(
			collection_route_view.Path_coordinates,
			domain.Route_path_coordinate{
				Stop_order:     route_stop_view_item.Stop_order,
				Bin_identifier: route_stop_view_item.Bin_identifier,
				Bin_code:       route_stop_view_item.Bin_code,
				Latitude:       route_stop_view_item.Latitude,
				Longitude:      route_stop_view_item.Longitude,
			},
		)
	}

	return collection_route_view, nil
}

func (route_repository *Route_repository) refresh_route_road_path_coordinates(
	application_context context.Context,
	query_executor postgres_query_executor,
	route_identifier string,
	route_stop_view_list []domain.Route_stop_view,
) error {
	ordered_route_stop_view_list := route_stop_view_list
	if ordered_route_stop_view_list == nil {
		loaded_route_stop_view_list, load_route_stops_error := load_route_stop_view_list_by_route_identifier(
			application_context,
			query_executor,
			route_identifier,
		)
		if load_route_stops_error != nil {
			return load_route_stops_error
		}
		ordered_route_stop_view_list = loaded_route_stop_view_list
	}

	ordered_path_coordinate_list := build_ordered_path_coordinates_from_route_stops(ordered_route_stop_view_list)
	computed_stops_hash := build_route_stops_hash(ordered_path_coordinate_list)

	load_existing_routing_state_statement := `
		SELECT
			COALESCE(stops_hash, ''),
			COALESCE(road_path_coordinates, '[]'::jsonb)
		FROM collection_routes
		WHERE id = $1::uuid
		FOR UPDATE;
	`
	var existing_stops_hash string
	var existing_raw_road_path_coordinates []byte
	load_existing_routing_state_error := query_executor.QueryRow(
		application_context,
		load_existing_routing_state_statement,
		route_identifier,
	).Scan(
		&existing_stops_hash,
		&existing_raw_road_path_coordinates,
	)
	if load_existing_routing_state_error != nil {
		return fmt.Errorf("failed_to_load_existing_route_routing_state: %w", load_existing_routing_state_error)
	}

	existing_has_road_path_coordinates := has_non_empty_road_path_coordinates(existing_raw_road_path_coordinates)
	if existing_stops_hash == computed_stops_hash && (existing_has_road_path_coordinates || len(ordered_path_coordinate_list) < 2) {
		return nil
	}

	resolved_road_path_coordinate_list := make([]domain.Route_road_path_coordinate, 0)
	routing_status := routing_status_not_required
	routing_error_text := ""
	if len(ordered_path_coordinate_list) >= 2 {
		if route_repository.route_path_router == nil {
			resolved_road_path_coordinate_list = build_straight_line_road_path_coordinates(
				ordered_path_coordinate_list,
			)
			routing_status = routing_status_fallback_straight_line
			routing_error_text = "route_path_router_not_configured"
		} else {
			calculated_road_path_coordinate_list, calculate_road_path_error := route_repository.route_path_router.calculate_road_path_coordinates(
				application_context,
				ordered_path_coordinate_list,
			)
			if calculate_road_path_error != nil {
				resolved_road_path_coordinate_list = build_straight_line_road_path_coordinates(
					ordered_path_coordinate_list,
				)
				routing_status = routing_status_fallback_straight_line
				routing_error_text = calculate_road_path_error.Error()
			} else {
				resolved_road_path_coordinate_list = calculated_road_path_coordinate_list
				routing_status = routing_status_routed
			}
		}
	}

	serialized_road_path_coordinate_list, marshal_road_path_error := json.Marshal(resolved_road_path_coordinate_list)
	if marshal_road_path_error != nil {
		return fmt.Errorf("failed_to_marshal_route_road_path_coordinates: %w", marshal_road_path_error)
	}

	update_route_routing_statement := `
		UPDATE collection_routes
		SET
			stops_hash = $2,
			road_path_coordinates = $3::jsonb,
			routed_at = NOW(),
			routing_provider = $4,
			routing_status = $5,
			routing_error = NULLIF($6, ''),
			updated_at = NOW()
		WHERE id = $1::uuid;
	`
	_, update_route_routing_error := query_executor.Exec(
		application_context,
		update_route_routing_statement,
		route_identifier,
		computed_stops_hash,
		serialized_road_path_coordinate_list,
		routing_provider_osrm,
		routing_status,
		strings.TrimSpace(routing_error_text),
	)
	if update_route_routing_error != nil {
		return fmt.Errorf("failed_to_update_route_routing_state: %w", update_route_routing_error)
	}

	return nil
}

func build_ordered_path_coordinates_from_route_stops(
	route_stop_view_list []domain.Route_stop_view,
) []domain.Route_path_coordinate {
	ordered_route_stop_view_list := append([]domain.Route_stop_view(nil), route_stop_view_list...)
	sort.SliceStable(
		ordered_route_stop_view_list,
		func(left_index int, right_index int) bool {
			return ordered_route_stop_view_list[left_index].Stop_order < ordered_route_stop_view_list[right_index].Stop_order
		},
	)

	ordered_path_coordinate_list := make(
		[]domain.Route_path_coordinate,
		0,
		len(ordered_route_stop_view_list),
	)
	for _, route_stop_view_item := range ordered_route_stop_view_list {
		ordered_path_coordinate_list = append(
			ordered_path_coordinate_list,
			domain.Route_path_coordinate{
				Stop_order:     route_stop_view_item.Stop_order,
				Bin_identifier: route_stop_view_item.Bin_identifier,
				Bin_code:       route_stop_view_item.Bin_code,
				Latitude:       route_stop_view_item.Latitude,
				Longitude:      route_stop_view_item.Longitude,
			},
		)
	}

	return ordered_path_coordinate_list
}

func build_route_stops_hash(ordered_path_coordinate_list []domain.Route_path_coordinate) string {
	hash_material_builder := strings.Builder{}
	for _, path_coordinate := range ordered_path_coordinate_list {
		hash_material_builder.WriteString(strconv.Itoa(path_coordinate.Stop_order))
		hash_material_builder.WriteString("|")
		hash_material_builder.WriteString(strconv.FormatFloat(path_coordinate.Latitude, 'f', 7, 64))
		hash_material_builder.WriteString(",")
		hash_material_builder.WriteString(strconv.FormatFloat(path_coordinate.Longitude, 'f', 7, 64))
		hash_material_builder.WriteString(";")
	}

	hashed_material := sha256.Sum256([]byte(hash_material_builder.String()))
	return hex.EncodeToString(hashed_material[:])
}

func build_straight_line_road_path_coordinates(
	ordered_path_coordinate_list []domain.Route_path_coordinate,
) []domain.Route_road_path_coordinate {
	straight_line_coordinate_list := make(
		[]domain.Route_road_path_coordinate,
		0,
		len(ordered_path_coordinate_list),
	)
	for _, path_coordinate := range ordered_path_coordinate_list {
		straight_line_coordinate_list = append(
			straight_line_coordinate_list,
			domain.Route_road_path_coordinate{
				path_coordinate.Longitude,
				path_coordinate.Latitude,
			},
		)
	}

	return straight_line_coordinate_list
}

func has_non_empty_road_path_coordinates(raw_road_path_coordinates []byte) bool {
	if len(raw_road_path_coordinates) == 0 {
		return false
	}

	road_path_coordinate_list := make([]domain.Route_road_path_coordinate, 0)
	unmarshal_road_path_error := json.Unmarshal(raw_road_path_coordinates, &road_path_coordinate_list)
	if unmarshal_road_path_error != nil {
		return false
	}

	return len(road_path_coordinate_list) > 0
}

func load_truck_route_assignment_by_identifier(
	application_context context.Context,
	query_executor postgres_query_executor,
	assignment_identifier string,
) (*domain.Truck_route_assignment_view, error) {
	load_assignment_statement := `
		SELECT
			assignment.id::text,
			assignment.truck_identifier,
			route.id::text,
			route.route_code,
			route.route_name,
			route.zone_name,
			route.collection_weekday::int,
			assignment.is_active,
			COALESCE(assignment.assigned_by_user_identifier, ''),
			COALESCE(assignment.assignment_notes, ''),
			assignment.assigned_at,
			assignment.unassigned_at
		FROM truck_route_assignments AS assignment
		INNER JOIN collection_routes AS route
			ON route.id = assignment.route_id
		WHERE assignment.id = $1::uuid;
	`

	assignment_view := &domain.Truck_route_assignment_view{}
	load_assignment_error := query_executor.QueryRow(
		application_context,
		load_assignment_statement,
		assignment_identifier,
	).Scan(
		&assignment_view.Assignment_identifier,
		&assignment_view.Truck_identifier,
		&assignment_view.Route_identifier,
		&assignment_view.Route_code,
		&assignment_view.Route_name,
		&assignment_view.Zone_name,
		&assignment_view.Collection_weekday,
		&assignment_view.Is_active,
		&assignment_view.Assigned_by_user_identifier,
		&assignment_view.Assignment_notes,
		&assignment_view.Assigned_at,
		&assignment_view.Unassigned_at,
	)
	if load_assignment_error != nil {
		if errors.Is(load_assignment_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("assignment_not_found")
		}
		return nil, fmt.Errorf("failed_to_load_truck_route_assignment_by_identifier: %w", load_assignment_error)
	}

	assignment_view.Assigned_by_user_identifier = strings.TrimSpace(assignment_view.Assigned_by_user_identifier)
	assignment_view.Assignment_notes = strings.TrimSpace(assignment_view.Assignment_notes)
	return assignment_view, nil
}

func resolve_active_route_identifier_by_truck_identifier(
	application_context context.Context,
	query_executor postgres_query_executor,
	truck_identifier string,
) (string, error) {
	resolve_assignment_statement := `
		SELECT route_id::text
		FROM truck_route_assignments
		WHERE truck_identifier = $1
			AND is_active = TRUE
		ORDER BY assigned_at DESC
		LIMIT 1;
	`

	var route_identifier string
	resolve_assignment_error := query_executor.QueryRow(
		application_context,
		resolve_assignment_statement,
		truck_identifier,
	).Scan(&route_identifier)
	if resolve_assignment_error != nil {
		if errors.Is(resolve_assignment_error, pgx.ErrNoRows) {
			return "", fmt.Errorf("active_route_assignment_not_found_for_truck")
		}
		return "", fmt.Errorf("failed_to_resolve_active_route_assignment_for_truck: %w", resolve_assignment_error)
	}

	return strings.TrimSpace(route_identifier), nil
}

func derive_route_deviation_severity_level(
	distance_to_route_meters float64,
	deviation_threshold_meters float64,
) string {
	if deviation_threshold_meters <= 0 {
		return "medium"
	}

	deviation_ratio := distance_to_route_meters / deviation_threshold_meters
	if deviation_ratio >= 2 {
		return "high"
	}
	if deviation_ratio >= 1.2 {
		return "medium"
	}
	return "low"
}

func normalize_route_deviation_alert_status_filter(raw_status_filter string) (string, error) {
	normalized_status_filter := strings.ToLower(strings.TrimSpace(raw_status_filter))
	switch normalized_status_filter {
	case "":
		return "", nil
	case "open", "resolved", "dismissed":
		return normalized_status_filter, nil
	default:
		return "", fmt.Errorf("invalid_route_deviation_alert_status_filter")
	}
}

func is_unique_constraint_error(database_error error, constraint_name string) bool {
	var postgres_error *pgconn.PgError
	if !errors.As(database_error, &postgres_error) {
		return false
	}

	return postgres_error.Code == "23505" && postgres_error.ConstraintName == constraint_name
}

func round_float_value(raw_float_value float64, decimal_places int) float64 {
	if decimal_places < 0 {
		return raw_float_value
	}

	decimal_multiplier := math.Pow10(decimal_places)
	return math.Round(raw_float_value*decimal_multiplier) / decimal_multiplier
}
