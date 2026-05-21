package storage

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"ecochitas/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Operations_repository struct {
	postgres_pool *pgxpool.Pool
}

var uuid_identifier_regex = regexp.MustCompile(
	`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func New_operations_repository(postgres_pool *pgxpool.Pool) *Operations_repository {
	return &Operations_repository{
		postgres_pool: postgres_pool,
	}
}

func (operations_repository *Operations_repository) Record_driver_collection_event(
	application_context context.Context,
	driver_collection_event_create_command domain.Driver_collection_event_create_command,
) (*domain.Driver_collection_event_record, error) {
	normalized_action_type, normalize_action_type_error := normalize_driver_action_type(
		driver_collection_event_create_command.Action_type,
	)
	if normalize_action_type_error != nil {
		return nil, normalize_action_type_error
	}

	if strings.TrimSpace(driver_collection_event_create_command.Bin_identifier) == "" {
		return nil, fmt.Errorf("bin_identifier_is_required")
	}

	if strings.TrimSpace(driver_collection_event_create_command.Authenticated_user_identifier) == "" {
		return nil, fmt.Errorf("authenticated_user_identifier_is_required")
	}

	if driver_collection_event_create_command.Action_at.IsZero() {
		driver_collection_event_create_command.Action_at = time.Now().UTC()
	}
	normalized_driver_user_identifier := strings.TrimSpace(
		driver_collection_event_create_command.Authenticated_user_identifier,
	)

	transaction_handler, begin_transaction_error := operations_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_driver_collection_event_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	driver_user_uuid_identifier := ""
	if is_valid_uuid_string(normalized_driver_user_identifier) {
		driver_user_uuid_identifier = normalized_driver_user_identifier

		ensure_driver_user_error := ensure_operation_user_exists(
			application_context,
			transaction_handler,
			normalized_driver_user_identifier,
			driver_collection_event_create_command.Authenticated_full_name,
			driver_collection_event_create_command.Authenticated_role_name,
		)
		if ensure_driver_user_error != nil {
			return nil, ensure_driver_user_error
		}
	}

	insert_event_statement := `
		INSERT INTO driver_collection_events (
			route_stop_id,
			bin_id,
			driver_user_id,
			driver_user_identifier,
			action_type,
			evidence_photo_url,
			action_notes,
			action_at
		) VALUES (
			NULLIF($1, '')::uuid,
			$2::uuid,
			NULLIF($3, '')::uuid,
			$4,
			$5,
			NULLIF($6, ''),
			NULLIF($7, ''),
			$8
		)
		RETURNING id::text;
	`

	var created_event_identifier string
	insert_event_error := transaction_handler.QueryRow(
		application_context,
		insert_event_statement,
		strings.TrimSpace(driver_collection_event_create_command.Route_stop_identifier),
		driver_collection_event_create_command.Bin_identifier,
		driver_user_uuid_identifier,
		normalized_driver_user_identifier,
		normalized_action_type,
		strings.TrimSpace(driver_collection_event_create_command.Evidence_photo_url),
		strings.TrimSpace(driver_collection_event_create_command.Action_notes),
		driver_collection_event_create_command.Action_at.UTC(),
	).Scan(&created_event_identifier)
	if insert_event_error != nil {
		return nil, fmt.Errorf("failed_to_insert_driver_collection_event: %w", insert_event_error)
	}

	update_bin_statement := `
		UPDATE bins
		SET
			fill_percentage = CASE
				WHEN $2 = 'emptied' THEN 0
				ELSE fill_percentage
			END,
			bin_status = CASE
				WHEN $2 = 'emptied' THEN 'available'
				WHEN $2 = 'contaminated' THEN 'warning'
				ELSE bin_status
			END,
			last_emptied_at = CASE
				WHEN $2 = 'emptied' THEN $3
				ELSE last_emptied_at
			END,
			updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING
			id::text,
			fill_percentage,
			bin_status,
			sensor_status,
			last_emptied_at;
	`

	driver_collection_event_record := &domain.Driver_collection_event_record{
		Event_identifier:       created_event_identifier,
		Bin_identifier:         driver_collection_event_create_command.Bin_identifier,
		Route_stop_identifier:  strings.TrimSpace(driver_collection_event_create_command.Route_stop_identifier),
		Driver_user_identifier: normalized_driver_user_identifier,
		Action_type:            normalized_action_type,
		Evidence_photo_url:     strings.TrimSpace(driver_collection_event_create_command.Evidence_photo_url),
		Action_notes:           strings.TrimSpace(driver_collection_event_create_command.Action_notes),
		Action_at:              driver_collection_event_create_command.Action_at.UTC(),
	}
	update_bin_error := transaction_handler.QueryRow(
		application_context,
		update_bin_statement,
		driver_collection_event_create_command.Bin_identifier,
		normalized_action_type,
		driver_collection_event_create_command.Action_at.UTC(),
	).Scan(
		&driver_collection_event_record.Bin_identifier,
		&driver_collection_event_record.Bin_fill_percentage,
		&driver_collection_event_record.Bin_status,
		&driver_collection_event_record.Bin_sensor_status,
		&driver_collection_event_record.Bin_last_emptied_at,
	)
	if update_bin_error != nil {
		if errors.Is(update_bin_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("bin_not_found")
		}
		return nil, fmt.Errorf("failed_to_update_bin_after_driver_event: %w", update_bin_error)
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_driver_collection_event_transaction: %w", commit_transaction_error)
	}

	return driver_collection_event_record, nil
}

func (operations_repository *Operations_repository) Ingest_bin_sensor_event(
	application_context context.Context,
	bin_sensor_event_ingestion_command domain.Bin_sensor_event_ingestion_command,
) (*domain.Bin_sensor_event_record, error) {
	normalized_sensor_status, normalize_sensor_status_error := normalize_sensor_status(
		bin_sensor_event_ingestion_command.Sensor_status,
	)
	if normalize_sensor_status_error != nil {
		return nil, normalize_sensor_status_error
	}

	if strings.TrimSpace(bin_sensor_event_ingestion_command.Bin_identifier) == "" {
		return nil, fmt.Errorf("bin_identifier_is_required")
	}

	if bin_sensor_event_ingestion_command.Fill_percentage < 0 || bin_sensor_event_ingestion_command.Fill_percentage > 100 {
		return nil, fmt.Errorf("fill_percentage_out_of_range")
	}

	if bin_sensor_event_ingestion_command.Measured_at.IsZero() {
		bin_sensor_event_ingestion_command.Measured_at = time.Now().UTC()
	}

	derived_bin_status := derive_bin_status_from_fill_percentage(bin_sensor_event_ingestion_command.Fill_percentage)
	if normalized_sensor_status == "offline" {
		derived_bin_status = "unknown"
	}

	transaction_handler, begin_transaction_error := operations_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_bin_sensor_event_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	insert_sensor_event_statement := `
		INSERT INTO bin_sensor_events (
			bin_id,
			fill_percentage,
			sensor_status,
			measured_at,
			source_identifier
		) VALUES (
			$1::uuid,
			$2,
			$3,
			$4,
			NULLIF($5, '')
		)
		RETURNING id::text, created_at;
	`

	bin_sensor_event_record := &domain.Bin_sensor_event_record{
		Bin_identifier:    strings.TrimSpace(bin_sensor_event_ingestion_command.Bin_identifier),
		Source_identifier: strings.TrimSpace(bin_sensor_event_ingestion_command.Source_identifier),
		Fill_percentage:   bin_sensor_event_ingestion_command.Fill_percentage,
		Sensor_status:     normalized_sensor_status,
		Bin_status:        derived_bin_status,
		Measured_at:       bin_sensor_event_ingestion_command.Measured_at.UTC(),
	}
	insert_sensor_event_error := transaction_handler.QueryRow(
		application_context,
		insert_sensor_event_statement,
		bin_sensor_event_ingestion_command.Bin_identifier,
		bin_sensor_event_ingestion_command.Fill_percentage,
		normalized_sensor_status,
		bin_sensor_event_ingestion_command.Measured_at.UTC(),
		bin_sensor_event_ingestion_command.Source_identifier,
	).Scan(
		&bin_sensor_event_record.Sensor_event_identifier,
		&bin_sensor_event_record.Recorded_at,
	)
	if insert_sensor_event_error != nil {
		return nil, fmt.Errorf("failed_to_insert_bin_sensor_event: %w", insert_sensor_event_error)
	}

	update_bin_statement := `
		UPDATE bins
		SET
			fill_percentage = $2,
			sensor_status = $3,
			bin_status = $4,
			updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING last_emptied_at;
	`
	update_bin_error := transaction_handler.QueryRow(
		application_context,
		update_bin_statement,
		bin_sensor_event_ingestion_command.Bin_identifier,
		bin_sensor_event_ingestion_command.Fill_percentage,
		normalized_sensor_status,
		derived_bin_status,
	).Scan(&bin_sensor_event_record.Bin_last_emptied_at)
	if update_bin_error != nil {
		if errors.Is(update_bin_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("bin_not_found")
		}
		return nil, fmt.Errorf("failed_to_update_bin_from_sensor_event: %w", update_bin_error)
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_bin_sensor_event_transaction: %w", commit_transaction_error)
	}

	return bin_sensor_event_record, nil
}

func (operations_repository *Operations_repository) Create_route_blockage_report(
	application_context context.Context,
	route_blockage_report_create_command domain.Route_blockage_report_create_command,
) (*domain.Route_blockage_report_record, error) {
	normalized_severity_level, normalize_severity_level_error := normalize_blockage_severity_level(
		route_blockage_report_create_command.Severity_level,
	)
	if normalize_severity_level_error != nil {
		return nil, normalize_severity_level_error
	}

	if strings.TrimSpace(route_blockage_report_create_command.Authenticated_user_identifier) == "" {
		return nil, fmt.Errorf("authenticated_user_identifier_is_required")
	}

	if strings.TrimSpace(route_blockage_report_create_command.Blockage_reason) == "" {
		return nil, fmt.Errorf("blockage_reason_is_required")
	}

	if route_blockage_report_create_command.Reported_at.IsZero() {
		route_blockage_report_create_command.Reported_at = time.Now().UTC()
	}

	insert_blockage_report_statement := `
		INSERT INTO collection_route_blockage_events (
			route_id,
			route_stop_id,
			bin_id,
			reported_by_user_identifier,
			blockage_reason,
			evidence_photo_url,
			severity_level,
			reported_at
		) VALUES (
			NULLIF($1, '')::uuid,
			NULLIF($2, '')::uuid,
			NULLIF($3, '')::uuid,
			$4,
			$5,
			NULLIF($6, ''),
			$7,
			$8
		)
		RETURNING
			id::text,
			COALESCE(route_id::text, ''),
			COALESCE(route_stop_id::text, ''),
			COALESCE(bin_id::text, ''),
			reported_by_user_identifier,
			blockage_reason,
			COALESCE(evidence_photo_url, ''),
			severity_level,
			status,
			reported_at,
			resolved_at,
			COALESCE(resolution_notes, ''),
			created_at,
			updated_at;
	`

	route_blockage_report_record := &domain.Route_blockage_report_record{}
	insert_blockage_report_error := operations_repository.postgres_pool.QueryRow(
		application_context,
		insert_blockage_report_statement,
		strings.TrimSpace(route_blockage_report_create_command.Route_identifier),
		strings.TrimSpace(route_blockage_report_create_command.Route_stop_identifier),
		strings.TrimSpace(route_blockage_report_create_command.Bin_identifier),
		strings.TrimSpace(route_blockage_report_create_command.Authenticated_user_identifier),
		strings.TrimSpace(route_blockage_report_create_command.Blockage_reason),
		strings.TrimSpace(route_blockage_report_create_command.Evidence_photo_url),
		normalized_severity_level,
		route_blockage_report_create_command.Reported_at.UTC(),
	).Scan(
		&route_blockage_report_record.Blockage_identifier,
		&route_blockage_report_record.Route_identifier,
		&route_blockage_report_record.Route_stop_identifier,
		&route_blockage_report_record.Bin_identifier,
		&route_blockage_report_record.Reported_by_user_identifier,
		&route_blockage_report_record.Blockage_reason,
		&route_blockage_report_record.Evidence_photo_url,
		&route_blockage_report_record.Severity_level,
		&route_blockage_report_record.Status,
		&route_blockage_report_record.Reported_at,
		&route_blockage_report_record.Resolved_at,
		&route_blockage_report_record.Resolution_notes,
		&route_blockage_report_record.Created_at,
		&route_blockage_report_record.Updated_at,
	)
	if insert_blockage_report_error != nil {
		return nil, fmt.Errorf("failed_to_insert_route_blockage_report: %w", insert_blockage_report_error)
	}

	route_blockage_report_record.Route_identifier = strings.TrimSpace(route_blockage_report_record.Route_identifier)
	route_blockage_report_record.Route_stop_identifier = strings.TrimSpace(route_blockage_report_record.Route_stop_identifier)
	route_blockage_report_record.Bin_identifier = strings.TrimSpace(route_blockage_report_record.Bin_identifier)
	route_blockage_report_record.Evidence_photo_url = strings.TrimSpace(route_blockage_report_record.Evidence_photo_url)
	route_blockage_report_record.Resolution_notes = strings.TrimSpace(route_blockage_report_record.Resolution_notes)

	return route_blockage_report_record, nil
}

func (operations_repository *Operations_repository) List_route_blockage_reports(
	application_context context.Context,
	route_blockage_report_list_query domain.Route_blockage_report_list_query,
) (*domain.Route_blockage_report_list_result, error) {
	normalized_status_filter, normalize_status_filter_error := normalize_blockage_status_filter(
		route_blockage_report_list_query.Status_filter,
	)
	if normalize_status_filter_error != nil {
		return nil, normalize_status_filter_error
	}

	normalized_limit := route_blockage_report_list_query.Limit
	if normalized_limit <= 0 {
		normalized_limit = 50
	}
	if normalized_limit > 200 {
		normalized_limit = 200
	}

	list_blockages_statement := `
		SELECT
			id::text,
			COALESCE(route_id::text, ''),
			COALESCE(route_stop_id::text, ''),
			COALESCE(bin_id::text, ''),
			reported_by_user_identifier,
			blockage_reason,
			COALESCE(evidence_photo_url, ''),
			severity_level,
			status,
			reported_at,
			resolved_at,
			COALESCE(resolution_notes, ''),
			created_at,
			updated_at
		FROM collection_route_blockage_events
		WHERE
			($1 = '' OR route_id = $1::uuid) AND
			($2 = '' OR route_stop_id = $2::uuid) AND
			($3 = '' OR bin_id = $3::uuid) AND
			($4 = '' OR status = $4)
		ORDER BY reported_at DESC, created_at DESC
		LIMIT $5;
	`

	query_rows, query_rows_error := operations_repository.postgres_pool.Query(
		application_context,
		list_blockages_statement,
		strings.TrimSpace(route_blockage_report_list_query.Route_identifier),
		strings.TrimSpace(route_blockage_report_list_query.Route_stop_identifier),
		strings.TrimSpace(route_blockage_report_list_query.Bin_identifier),
		normalized_status_filter,
		normalized_limit,
	)
	if query_rows_error != nil {
		return nil, fmt.Errorf("failed_to_query_route_blockage_reports: %w", query_rows_error)
	}
	defer query_rows.Close()

	blockage_report_item_list := make([]domain.Route_blockage_report_record, 0)
	for query_rows.Next() {
		var blockage_report_item domain.Route_blockage_report_record
		scan_row_error := query_rows.Scan(
			&blockage_report_item.Blockage_identifier,
			&blockage_report_item.Route_identifier,
			&blockage_report_item.Route_stop_identifier,
			&blockage_report_item.Bin_identifier,
			&blockage_report_item.Reported_by_user_identifier,
			&blockage_report_item.Blockage_reason,
			&blockage_report_item.Evidence_photo_url,
			&blockage_report_item.Severity_level,
			&blockage_report_item.Status,
			&blockage_report_item.Reported_at,
			&blockage_report_item.Resolved_at,
			&blockage_report_item.Resolution_notes,
			&blockage_report_item.Created_at,
			&blockage_report_item.Updated_at,
		)
		if scan_row_error != nil {
			return nil, fmt.Errorf("failed_to_scan_route_blockage_report_row: %w", scan_row_error)
		}

		blockage_report_item.Route_identifier = strings.TrimSpace(blockage_report_item.Route_identifier)
		blockage_report_item.Route_stop_identifier = strings.TrimSpace(blockage_report_item.Route_stop_identifier)
		blockage_report_item.Bin_identifier = strings.TrimSpace(blockage_report_item.Bin_identifier)
		blockage_report_item.Evidence_photo_url = strings.TrimSpace(blockage_report_item.Evidence_photo_url)
		blockage_report_item.Resolution_notes = strings.TrimSpace(blockage_report_item.Resolution_notes)

		blockage_report_item_list = append(blockage_report_item_list, blockage_report_item)
	}

	iterate_rows_error := query_rows.Err()
	if iterate_rows_error != nil {
		return nil, fmt.Errorf("failed_during_route_blockage_reports_rows_iteration: %w", iterate_rows_error)
	}

	return &domain.Route_blockage_report_list_result{
		Items: blockage_report_item_list,
		Total: len(blockage_report_item_list),
	}, nil
}

func (operations_repository *Operations_repository) Update_route_blockage_report_status(
	application_context context.Context,
	route_blockage_report_status_update_command domain.Route_blockage_report_status_update_command,
) (*domain.Route_blockage_report_record, error) {
	normalized_blockage_identifier := strings.TrimSpace(
		route_blockage_report_status_update_command.Blockage_identifier,
	)
	if !is_valid_uuid_string(normalized_blockage_identifier) {
		return nil, fmt.Errorf("invalid_blockage_identifier")
	}

	normalized_status, normalize_status_error := normalize_blockage_resolution_status(
		route_blockage_report_status_update_command.Status,
	)
	if normalize_status_error != nil {
		return nil, normalize_status_error
	}

	resolved_at_value := route_blockage_report_status_update_command.Resolved_at.UTC()
	if route_blockage_report_status_update_command.Resolved_at.IsZero() {
		resolved_at_value = time.Now().UTC()
	}

	update_blockage_statement := `
		UPDATE collection_route_blockage_events
		SET
			status = $2,
			resolved_at = $3,
			resolution_notes = NULLIF($4, ''),
			updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING
			id::text,
			COALESCE(route_id::text, ''),
			COALESCE(route_stop_id::text, ''),
			COALESCE(bin_id::text, ''),
			reported_by_user_identifier,
			blockage_reason,
			COALESCE(evidence_photo_url, ''),
			severity_level,
			status,
			reported_at,
			resolved_at,
			COALESCE(resolution_notes, ''),
			created_at,
			updated_at;
	`

	updated_route_blockage_report_record := &domain.Route_blockage_report_record{}
	update_blockage_error := operations_repository.postgres_pool.QueryRow(
		application_context,
		update_blockage_statement,
		normalized_blockage_identifier,
		normalized_status,
		resolved_at_value,
		strings.TrimSpace(route_blockage_report_status_update_command.Resolution_notes),
	).Scan(
		&updated_route_blockage_report_record.Blockage_identifier,
		&updated_route_blockage_report_record.Route_identifier,
		&updated_route_blockage_report_record.Route_stop_identifier,
		&updated_route_blockage_report_record.Bin_identifier,
		&updated_route_blockage_report_record.Reported_by_user_identifier,
		&updated_route_blockage_report_record.Blockage_reason,
		&updated_route_blockage_report_record.Evidence_photo_url,
		&updated_route_blockage_report_record.Severity_level,
		&updated_route_blockage_report_record.Status,
		&updated_route_blockage_report_record.Reported_at,
		&updated_route_blockage_report_record.Resolved_at,
		&updated_route_blockage_report_record.Resolution_notes,
		&updated_route_blockage_report_record.Created_at,
		&updated_route_blockage_report_record.Updated_at,
	)
	if update_blockage_error != nil {
		if errors.Is(update_blockage_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("blockage_report_not_found")
		}
		return nil, fmt.Errorf("failed_to_update_route_blockage_report_status: %w", update_blockage_error)
	}

	updated_route_blockage_report_record.Route_identifier = strings.TrimSpace(updated_route_blockage_report_record.Route_identifier)
	updated_route_blockage_report_record.Route_stop_identifier = strings.TrimSpace(updated_route_blockage_report_record.Route_stop_identifier)
	updated_route_blockage_report_record.Bin_identifier = strings.TrimSpace(updated_route_blockage_report_record.Bin_identifier)
	updated_route_blockage_report_record.Evidence_photo_url = strings.TrimSpace(updated_route_blockage_report_record.Evidence_photo_url)
	updated_route_blockage_report_record.Resolution_notes = strings.TrimSpace(updated_route_blockage_report_record.Resolution_notes)

	return updated_route_blockage_report_record, nil
}

func ensure_operation_user_exists(
	application_context context.Context,
	query_executor postgres_query_executor,
	operation_user_identifier string,
	operation_user_full_name string,
	operation_user_role_name string,
) error {
	normalized_user_role_name := strings.ToLower(strings.TrimSpace(operation_user_role_name))
	switch normalized_user_role_name {
	case "driver", "admin", "citizen", "condominium_admin":
	default:
		return fmt.Errorf("invalid_operation_user_role_name")
	}

	upsert_user_statement := `
		INSERT INTO users (
			id,
			full_name,
			role_name,
			zone_name,
			home_reference
		) VALUES (
			$1::uuid,
			COALESCE(NULLIF($2, ''), 'Unknown User'),
			$3,
			NULL,
			NULL
		)
		ON CONFLICT (id)
		DO UPDATE SET
			full_name = COALESCE(NULLIF(EXCLUDED.full_name, ''), users.full_name),
			role_name = EXCLUDED.role_name,
			updated_at = NOW();
	`

	_, upsert_user_error := query_executor.Exec(
		application_context,
		upsert_user_statement,
		strings.TrimSpace(operation_user_identifier),
		strings.TrimSpace(operation_user_full_name),
		normalized_user_role_name,
	)
	if upsert_user_error != nil {
		return fmt.Errorf("failed_to_upsert_operation_user: %w", upsert_user_error)
	}

	return nil
}

func normalize_driver_action_type(raw_driver_action_type string) (string, error) {
	normalized_driver_action_type := strings.ToLower(strings.TrimSpace(raw_driver_action_type))
	switch normalized_driver_action_type {
	case "emptied", "not_accessible", "contaminated":
		return normalized_driver_action_type, nil
	default:
		return "", fmt.Errorf("invalid_action_type")
	}
}

func normalize_sensor_status(raw_sensor_status string) (string, error) {
	normalized_sensor_status := strings.ToLower(strings.TrimSpace(raw_sensor_status))
	switch normalized_sensor_status {
	case "online", "offline":
		return normalized_sensor_status, nil
	default:
		return "", fmt.Errorf("invalid_sensor_status")
	}
}

func normalize_blockage_severity_level(raw_severity_level string) (string, error) {
	normalized_severity_level := strings.ToLower(strings.TrimSpace(raw_severity_level))
	switch normalized_severity_level {
	case "", "medium":
		return "medium", nil
	case "low", "high":
		return normalized_severity_level, nil
	default:
		return "", fmt.Errorf("invalid_severity_level")
	}
}

func normalize_blockage_status_filter(raw_status_filter string) (string, error) {
	normalized_status_filter := strings.ToLower(strings.TrimSpace(raw_status_filter))
	switch normalized_status_filter {
	case "":
		return "", nil
	case "open", "resolved", "dismissed":
		return normalized_status_filter, nil
	default:
		return "", fmt.Errorf("invalid_status_filter")
	}
}

func normalize_blockage_resolution_status(raw_status string) (string, error) {
	normalized_status := strings.ToLower(strings.TrimSpace(raw_status))
	switch normalized_status {
	case "resolved", "dismissed":
		return normalized_status, nil
	default:
		return "", fmt.Errorf("invalid_blockage_status")
	}
}

func derive_bin_status_from_fill_percentage(fill_percentage int) string {
	if fill_percentage >= 90 {
		return "full"
	}
	if fill_percentage >= 70 {
		return "warning"
	}
	return "available"
}

func is_valid_uuid_string(candidate_uuid_string string) bool {
	if strings.TrimSpace(candidate_uuid_string) == "" {
		return false
	}

	return uuid_identifier_regex.MatchString(candidate_uuid_string)
}
