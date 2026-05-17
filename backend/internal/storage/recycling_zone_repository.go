package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"ecochitas/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres_query_executor interface {
	Exec(application_context context.Context, query_statement string, query_arguments ...any) (pgconn.CommandTag, error)
	Query(application_context context.Context, query_statement string, query_arguments ...any) (pgx.Rows, error)
	QueryRow(application_context context.Context, query_statement string, query_arguments ...any) pgx.Row
}

type Recycling_zone_repository struct {
	postgres_pool *pgxpool.Pool
}

func New_recycling_zone_repository(postgres_pool *pgxpool.Pool) *Recycling_zone_repository {
	return &Recycling_zone_repository{
		postgres_pool: postgres_pool,
	}
}

func (recycling_zone_repository *Recycling_zone_repository) List_zone_recycling_containers(
	application_context context.Context,
	zone_code_filter string,
) ([]domain.Zone_recycling_container_view, error) {
	query_statement := `
		SELECT
			container.id::text,
			container.container_code,
			zone.id::text,
			zone.zone_code,
			zone.zone_name,
			container.latitude::float8,
			container.longitude::float8,
			COUNT(DISTINCT assignment.household_id)::int AS assigned_household_total,
			COALESCE(eligible_household_subquery.eligible_household_total, 0)::int AS eligible_household_total,
			COALESCE(contaminated_cycle_subquery.contaminated_cycle_total, 0)::int AS contaminated_cycle_total,
			COALESCE(completed_cycle_subquery.completed_cycle_total, 0)::int AS completed_cycle_total
		FROM zone_recycling_containers AS container
		INNER JOIN recycling_zones AS zone
			ON zone.id = container.zone_id
		LEFT JOIN zone_recycling_container_household_assignments AS assignment
			ON assignment.container_id = container.id
			AND assignment.is_active = TRUE
		LEFT JOIN LATERAL (
			SELECT
				COUNT(DISTINCT evidence.household_id)::int AS eligible_household_total
			FROM recycling_collection_cycles AS cycle
			INNER JOIN recycling_cycle_evidence_submissions AS evidence
				ON evidence.cycle_id = cycle.id
			WHERE cycle.container_id = container.id
				AND cycle.cycle_status = 'closed'
				AND evidence.validation_status = 'accepted'
		) AS eligible_household_subquery ON TRUE
		LEFT JOIN LATERAL (
			SELECT COUNT(*)::int AS contaminated_cycle_total
			FROM recycling_collection_cycles AS cycle
			WHERE cycle.container_id = container.id
				AND cycle.cycle_status = 'closed'
				AND cycle.contamination_level IN ('medium', 'high')
		) AS contaminated_cycle_subquery ON TRUE
		LEFT JOIN LATERAL (
			SELECT COUNT(*)::int AS completed_cycle_total
			FROM recycling_collection_cycles AS cycle
			WHERE cycle.container_id = container.id
				AND cycle.cycle_status = 'closed'
		) AS completed_cycle_subquery ON TRUE
		WHERE container.is_active = TRUE
			AND zone.is_active = TRUE
			AND ($1 = '' OR zone.zone_code = $1)
		GROUP BY
			container.id,
			container.container_code,
			zone.id,
			zone.zone_code,
			zone.zone_name,
			container.latitude,
			container.longitude,
			eligible_household_subquery.eligible_household_total,
			contaminated_cycle_subquery.contaminated_cycle_total,
			completed_cycle_subquery.completed_cycle_total
		ORDER BY zone.zone_name, container.container_code;
	`

	query_rows, query_error := recycling_zone_repository.postgres_pool.Query(
		application_context,
		query_statement,
		strings.TrimSpace(zone_code_filter),
	)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_zone_recycling_containers: %w", query_error)
	}
	defer query_rows.Close()

	container_view_list := make([]domain.Zone_recycling_container_view, 0)
	for query_rows.Next() {
		var container_view_item domain.Zone_recycling_container_view
		scan_error := query_rows.Scan(
			&container_view_item.Container_identifier,
			&container_view_item.Container_code,
			&container_view_item.Zone_identifier,
			&container_view_item.Zone_code,
			&container_view_item.Zone_name,
			&container_view_item.Latitude,
			&container_view_item.Longitude,
			&container_view_item.Assigned_household_total,
			&container_view_item.Eligible_household_total,
			&container_view_item.Contaminated_cycle_total,
			&container_view_item.Completed_cycle_total,
		)
		if scan_error != nil {
			return nil, fmt.Errorf("failed_to_scan_zone_recycling_container_row: %w", scan_error)
		}

		container_view_list = append(container_view_list, container_view_item)
	}

	rows_iteration_error := query_rows.Err()
	if rows_iteration_error != nil {
		return nil, fmt.Errorf("failed_during_zone_recycling_containers_rows_iteration: %w", rows_iteration_error)
	}

	return container_view_list, nil
}

func (recycling_zone_repository *Recycling_zone_repository) Start_recycling_collection_cycle(
	application_context context.Context,
	recycling_cycle_start_command domain.Recycling_cycle_start_command,
) (*domain.Recycling_collection_cycle, error) {
	scheduled_collection_date := recycling_cycle_start_command.Scheduled_collection_date
	if scheduled_collection_date.IsZero() {
		scheduled_collection_date = time.Now().UTC()
	}

	insert_cycle_statement := `
		INSERT INTO recycling_collection_cycles (
			container_id,
			scheduled_collection_date,
			cycle_status,
			collection_operator_name,
			started_at
		) VALUES (
			$1::uuid,
			$2::date,
			'in_progress',
			NULLIF($3, ''),
			NOW()
		)
		RETURNING id::text;
	`

	var cycle_identifier string
	insert_cycle_error := recycling_zone_repository.postgres_pool.QueryRow(
		application_context,
		insert_cycle_statement,
		recycling_cycle_start_command.Container_identifier,
		scheduled_collection_date.Format("2006-01-02"),
		strings.TrimSpace(recycling_cycle_start_command.Collection_operator_name),
	).Scan(&cycle_identifier)
	if insert_cycle_error != nil {
		return nil, fmt.Errorf("failed_to_start_recycling_collection_cycle: %w", insert_cycle_error)
	}

	created_cycle, load_cycle_error := recycling_zone_repository.Get_recycling_collection_cycle(
		application_context,
		cycle_identifier,
	)
	if load_cycle_error != nil {
		return nil, load_cycle_error
	}

	return created_cycle, nil
}

func (recycling_zone_repository *Recycling_zone_repository) Submit_recycling_evidence(
	application_context context.Context,
	recycling_evidence_submission_command domain.Recycling_evidence_submission_command,
) (*domain.Recycling_evidence_submission, error) {
	normalized_validation_status := normalize_validation_status(recycling_evidence_submission_command.Validation_status)

	upsert_submission_statement := `
		WITH eligible_household_assignment AS (
			SELECT assignment.household_id
			FROM recycling_collection_cycles AS cycle
			INNER JOIN zone_recycling_container_household_assignments AS assignment
				ON assignment.container_id = cycle.container_id
				AND assignment.is_active = TRUE
			WHERE cycle.id = $1::uuid
				AND assignment.household_id = $2::uuid
		)
		INSERT INTO recycling_cycle_evidence_submissions (
			cycle_id,
			household_id,
			evidence_photo_url,
			evidence_captured_at,
			evidence_latitude,
			evidence_longitude,
			validation_status,
			rejection_reason
		)
		SELECT
			$1::uuid,
			$2::uuid,
			$3,
			$4,
			$5,
			$6,
			$7,
			NULLIF($8, '')
		FROM eligible_household_assignment
		ON CONFLICT (cycle_id, household_id)
		DO UPDATE SET
			evidence_photo_url = EXCLUDED.evidence_photo_url,
			evidence_captured_at = EXCLUDED.evidence_captured_at,
			evidence_latitude = EXCLUDED.evidence_latitude,
			evidence_longitude = EXCLUDED.evidence_longitude,
			validation_status = EXCLUDED.validation_status,
			rejection_reason = EXCLUDED.rejection_reason,
			updated_at = NOW()
		RETURNING
			id::text,
			cycle_id::text,
			household_id::text,
			validation_status,
			created_at,
			updated_at;
	`

	var recycling_evidence_submission domain.Recycling_evidence_submission
	upsert_submission_error := recycling_zone_repository.postgres_pool.QueryRow(
		application_context,
		upsert_submission_statement,
		recycling_evidence_submission_command.Cycle_identifier,
		recycling_evidence_submission_command.Household_identifier,
		recycling_evidence_submission_command.Evidence_photo_url,
		recycling_evidence_submission_command.Evidence_captured_at,
		recycling_evidence_submission_command.Evidence_latitude,
		recycling_evidence_submission_command.Evidence_longitude,
		normalized_validation_status,
		strings.TrimSpace(recycling_evidence_submission_command.Rejection_reason),
	).Scan(
		&recycling_evidence_submission.Submission_identifier,
		&recycling_evidence_submission.Cycle_identifier,
		&recycling_evidence_submission.Household_identifier,
		&recycling_evidence_submission.Validation_status,
		&recycling_evidence_submission.Created_at,
		&recycling_evidence_submission.Updated_at,
	)
	if upsert_submission_error != nil {
		if errors.Is(upsert_submission_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("household_is_not_assigned_to_cycle_container")
		}
		return nil, fmt.Errorf("failed_to_submit_recycling_evidence: %w", upsert_submission_error)
	}

	return &recycling_evidence_submission, nil
}

func (recycling_zone_repository *Recycling_zone_repository) Close_recycling_collection_cycle(
	application_context context.Context,
	recycling_cycle_close_command domain.Recycling_cycle_close_command,
) (*domain.Recycling_cycle_summary, error) {
	if strings.TrimSpace(recycling_cycle_close_command.Cycle_identifier) == "" {
		return nil, fmt.Errorf("cycle_identifier_is_required")
	}

	if recycling_cycle_close_command.Raw_points_total < 0 {
		return nil, fmt.Errorf("raw_points_total_must_be_greater_or_equal_than_zero")
	}

	normalized_contamination_level, normalize_contamination_level_error := normalize_contamination_level(
		recycling_cycle_close_command.Contamination_level,
	)
	if normalize_contamination_level_error != nil {
		return nil, normalize_contamination_level_error
	}

	contamination_discount_percentage := resolve_contamination_discount_percentage(
		normalized_contamination_level,
		recycling_cycle_close_command.Contamination_discount_percentage,
	)
	raw_points_total_cents := to_points_cents(recycling_cycle_close_command.Raw_points_total)
	discount_points_total_cents := int64(math.Round(float64(raw_points_total_cents) * contamination_discount_percentage / 100))
	final_points_total_cents := raw_points_total_cents - discount_points_total_cents
	if final_points_total_cents < 0 {
		final_points_total_cents = 0
	}

	transaction_handler, begin_transaction_error := recycling_zone_repository.postgres_pool.BeginTx(
		application_context,
		pgx.TxOptions{},
	)
	if begin_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_begin_close_cycle_transaction: %w", begin_transaction_error)
	}
	defer func() {
		_ = transaction_handler.Rollback(application_context)
	}()

	update_cycle_statement := `
		UPDATE recycling_collection_cycles
		SET
			cycle_status = 'closed',
			collection_operator_name = COALESCE(NULLIF($2, ''), collection_operator_name),
			raw_points_total = $3,
			contamination_level = $4,
			contamination_discount_percentage = $5,
			discount_points_total = $6,
			final_points_total = $7,
			contamination_notes = NULLIF($8, ''),
			closed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING id::text;
	`

	var updated_cycle_identifier string
	update_cycle_error := transaction_handler.QueryRow(
		application_context,
		update_cycle_statement,
		recycling_cycle_close_command.Cycle_identifier,
		strings.TrimSpace(recycling_cycle_close_command.Collection_operator_name),
		from_points_cents(raw_points_total_cents),
		normalized_contamination_level,
		contamination_discount_percentage,
		from_points_cents(discount_points_total_cents),
		from_points_cents(final_points_total_cents),
		strings.TrimSpace(recycling_cycle_close_command.Contamination_notes),
	).Scan(&updated_cycle_identifier)
	if update_cycle_error != nil {
		if errors.Is(update_cycle_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("recycling_collection_cycle_not_found")
		}
		return nil, fmt.Errorf("failed_to_update_recycling_collection_cycle: %w", update_cycle_error)
	}

	delete_previous_points_statement := `
		DELETE FROM recycling_cycle_household_points
		WHERE cycle_id = $1::uuid;
	`
	_, delete_previous_points_error := transaction_handler.Exec(
		application_context,
		delete_previous_points_statement,
		updated_cycle_identifier,
	)
	if delete_previous_points_error != nil {
		return nil, fmt.Errorf("failed_to_delete_previous_household_points: %w", delete_previous_points_error)
	}

	eligible_household_identifiers, load_eligible_households_error := load_eligible_household_identifiers(
		application_context,
		transaction_handler,
		updated_cycle_identifier,
	)
	if load_eligible_households_error != nil {
		return nil, load_eligible_households_error
	}

	insert_household_points_statement := `
		INSERT INTO recycling_cycle_household_points (
			cycle_id,
			household_id,
			awarded_points,
			updated_at
		) VALUES (
			$1::uuid,
			$2::uuid,
			$3,
			NOW()
		);
	`

	if len(eligible_household_identifiers) > 0 && final_points_total_cents > 0 {
		base_household_points_cents := final_points_total_cents / int64(len(eligible_household_identifiers))
		remainder_points_cents := final_points_total_cents % int64(len(eligible_household_identifiers))

		for eligible_household_index, eligible_household_identifier := range eligible_household_identifiers {
			awarded_points_cents := base_household_points_cents
			if int64(eligible_household_index) < remainder_points_cents {
				awarded_points_cents += 1
			}

			_, insert_household_points_error := transaction_handler.Exec(
				application_context,
				insert_household_points_statement,
				updated_cycle_identifier,
				eligible_household_identifier,
				from_points_cents(awarded_points_cents),
			)
			if insert_household_points_error != nil {
				return nil, fmt.Errorf("failed_to_insert_household_points: %w", insert_household_points_error)
			}
		}
	}

	commit_transaction_error := transaction_handler.Commit(application_context)
	if commit_transaction_error != nil {
		return nil, fmt.Errorf("failed_to_commit_close_cycle_transaction: %w", commit_transaction_error)
	}

	closed_cycle_summary, load_summary_error := recycling_zone_repository.Get_recycling_cycle_summary(
		application_context,
		updated_cycle_identifier,
	)
	if load_summary_error != nil {
		return nil, load_summary_error
	}

	return closed_cycle_summary, nil
}

func (recycling_zone_repository *Recycling_zone_repository) Get_recycling_collection_cycle(
	application_context context.Context,
	cycle_identifier string,
) (*domain.Recycling_collection_cycle, error) {
	query_statement := `
		SELECT
			cycle.id::text,
			container.id::text,
			container.container_code,
			zone.id::text,
			zone.zone_code,
			zone.zone_name,
			cycle.scheduled_collection_date::text,
			cycle.cycle_status,
			COALESCE(cycle.collection_operator_name, ''),
			cycle.raw_points_total::float8,
			cycle.contamination_level,
			cycle.contamination_discount_percentage::float8,
			cycle.discount_points_total::float8,
			cycle.final_points_total::float8,
			COALESCE(cycle.contamination_notes, ''),
			cycle.closed_at
		FROM recycling_collection_cycles AS cycle
		INNER JOIN zone_recycling_containers AS container
			ON container.id = cycle.container_id
		INNER JOIN recycling_zones AS zone
			ON zone.id = container.zone_id
		WHERE cycle.id = $1::uuid;
	`

	var recycling_collection_cycle domain.Recycling_collection_cycle
	query_cycle_error := recycling_zone_repository.postgres_pool.QueryRow(
		application_context,
		query_statement,
		cycle_identifier,
	).Scan(
		&recycling_collection_cycle.Cycle_identifier,
		&recycling_collection_cycle.Container_identifier,
		&recycling_collection_cycle.Container_code,
		&recycling_collection_cycle.Zone_identifier,
		&recycling_collection_cycle.Zone_code,
		&recycling_collection_cycle.Zone_name,
		&recycling_collection_cycle.Scheduled_collection_date,
		&recycling_collection_cycle.Cycle_status,
		&recycling_collection_cycle.Collection_operator_name,
		&recycling_collection_cycle.Raw_points_total,
		&recycling_collection_cycle.Contamination_level,
		&recycling_collection_cycle.Contamination_discount_percent,
		&recycling_collection_cycle.Discount_points_total,
		&recycling_collection_cycle.Final_points_total,
		&recycling_collection_cycle.Contamination_notes,
		&recycling_collection_cycle.Closed_at,
	)
	if query_cycle_error != nil {
		if errors.Is(query_cycle_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("recycling_collection_cycle_not_found")
		}
		return nil, fmt.Errorf("failed_to_query_recycling_collection_cycle: %w", query_cycle_error)
	}

	return &recycling_collection_cycle, nil
}

func (recycling_zone_repository *Recycling_zone_repository) Get_recycling_cycle_summary(
	application_context context.Context,
	cycle_identifier string,
) (*domain.Recycling_cycle_summary, error) {
	return load_recycling_cycle_summary(
		application_context,
		recycling_zone_repository.postgres_pool,
		cycle_identifier,
	)
}

func load_recycling_cycle_summary(
	application_context context.Context,
	query_executor postgres_query_executor,
	cycle_identifier string,
) (*domain.Recycling_cycle_summary, error) {
	summary_query_statement := `
		SELECT
			cycle.id::text,
			container.id::text,
			container.container_code,
			zone.id::text,
			zone.zone_code,
			zone.zone_name,
			cycle.cycle_status,
			cycle.scheduled_collection_date::text,
			COALESCE(cycle.collection_operator_name, ''),
			COALESCE(eligible_household_subquery.eligible_household_total, 0)::int,
			cycle.raw_points_total::float8,
			cycle.contamination_level,
			cycle.contamination_discount_percentage::float8,
			cycle.discount_points_total::float8,
			cycle.final_points_total::float8,
			COALESCE(cycle.contamination_notes, ''),
			cycle.closed_at
		FROM recycling_collection_cycles AS cycle
		INNER JOIN zone_recycling_containers AS container
			ON container.id = cycle.container_id
		INNER JOIN recycling_zones AS zone
			ON zone.id = container.zone_id
		LEFT JOIN LATERAL (
			SELECT COUNT(DISTINCT evidence.household_id)::int AS eligible_household_total
			FROM recycling_cycle_evidence_submissions AS evidence
			WHERE evidence.cycle_id = cycle.id
				AND evidence.validation_status = 'accepted'
		) AS eligible_household_subquery ON TRUE
		WHERE cycle.id = $1::uuid;
	`

	var recycling_cycle_summary domain.Recycling_cycle_summary
	load_summary_error := query_executor.QueryRow(
		application_context,
		summary_query_statement,
		cycle_identifier,
	).Scan(
		&recycling_cycle_summary.Cycle_identifier,
		&recycling_cycle_summary.Container_identifier,
		&recycling_cycle_summary.Container_code,
		&recycling_cycle_summary.Zone_identifier,
		&recycling_cycle_summary.Zone_code,
		&recycling_cycle_summary.Zone_name,
		&recycling_cycle_summary.Cycle_status,
		&recycling_cycle_summary.Scheduled_collection_date,
		&recycling_cycle_summary.Collection_operator_name,
		&recycling_cycle_summary.Eligible_household_total,
		&recycling_cycle_summary.Raw_points_total,
		&recycling_cycle_summary.Contamination_level,
		&recycling_cycle_summary.Contamination_discount_percentage,
		&recycling_cycle_summary.Discount_points_total,
		&recycling_cycle_summary.Final_points_total,
		&recycling_cycle_summary.Contamination_notes,
		&recycling_cycle_summary.Closed_at,
	)
	if load_summary_error != nil {
		if errors.Is(load_summary_error, pgx.ErrNoRows) {
			return nil, fmt.Errorf("recycling_collection_cycle_not_found")
		}
		return nil, fmt.Errorf("failed_to_query_recycling_cycle_summary: %w", load_summary_error)
	}

	household_points_query_statement := `
		SELECT
			household.id::text,
			household.household_code,
			household_points.awarded_points::float8
		FROM recycling_cycle_household_points AS household_points
		INNER JOIN zone_households AS household
			ON household.id = household_points.household_id
		WHERE household_points.cycle_id = $1::uuid
		ORDER BY household.household_code;
	`
	query_rows, household_points_query_error := query_executor.Query(
		application_context,
		household_points_query_statement,
		cycle_identifier,
	)
	if household_points_query_error != nil {
		return nil, fmt.Errorf("failed_to_query_household_points_summary: %w", household_points_query_error)
	}
	defer query_rows.Close()

	household_points_list := make([]domain.Recycling_cycle_household_points, 0)
	for query_rows.Next() {
		var household_points_item domain.Recycling_cycle_household_points
		scan_household_points_error := query_rows.Scan(
			&household_points_item.Household_identifier,
			&household_points_item.Household_code,
			&household_points_item.Awarded_points,
		)
		if scan_household_points_error != nil {
			return nil, fmt.Errorf("failed_to_scan_household_points_row: %w", scan_household_points_error)
		}

		household_points_list = append(household_points_list, household_points_item)
	}

	rows_iteration_error := query_rows.Err()
	if rows_iteration_error != nil {
		return nil, fmt.Errorf("failed_during_household_points_rows_iteration: %w", rows_iteration_error)
	}

	recycling_cycle_summary.Household_points = household_points_list
	return &recycling_cycle_summary, nil
}

func load_eligible_household_identifiers(
	application_context context.Context,
	query_executor postgres_query_executor,
	cycle_identifier string,
) ([]string, error) {
	query_statement := `
		SELECT DISTINCT evidence.household_id::text
		FROM recycling_cycle_evidence_submissions AS evidence
		WHERE evidence.cycle_id = $1::uuid
			AND evidence.validation_status = 'accepted'
		ORDER BY evidence.household_id::text;
	`

	query_rows, query_error := query_executor.Query(application_context, query_statement, cycle_identifier)
	if query_error != nil {
		return nil, fmt.Errorf("failed_to_query_eligible_households: %w", query_error)
	}
	defer query_rows.Close()

	eligible_household_identifier_list := make([]string, 0)
	for query_rows.Next() {
		var eligible_household_identifier string
		scan_error := query_rows.Scan(&eligible_household_identifier)
		if scan_error != nil {
			return nil, fmt.Errorf("failed_to_scan_eligible_household_identifier: %w", scan_error)
		}

		eligible_household_identifier_list = append(eligible_household_identifier_list, eligible_household_identifier)
	}

	rows_iteration_error := query_rows.Err()
	if rows_iteration_error != nil {
		return nil, fmt.Errorf("failed_during_eligible_households_rows_iteration: %w", rows_iteration_error)
	}

	return eligible_household_identifier_list, nil
}

func normalize_validation_status(raw_validation_status string) string {
	switch strings.ToLower(strings.TrimSpace(raw_validation_status)) {
	case "pending":
		return "pending"
	case "rejected":
		return "rejected"
	default:
		return "accepted"
	}
}

func normalize_contamination_level(raw_contamination_level string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw_contamination_level)) {
	case "", "low":
		return "low", nil
	case "medium":
		return "medium", nil
	case "high":
		return "high", nil
	default:
		return "", fmt.Errorf("invalid_contamination_level")
	}
}

func resolve_contamination_discount_percentage(contamination_level string, custom_discount_percentage *float64) float64 {
	if custom_discount_percentage != nil {
		clamped_discount_percentage := math.Max(0, math.Min(100, *custom_discount_percentage))
		return clamped_discount_percentage
	}

	switch contamination_level {
	case "medium":
		return 10
	case "high":
		return 25
	default:
		return 0
	}
}

func to_points_cents(points_amount float64) int64 {
	return int64(math.Round(points_amount * 100))
}

func from_points_cents(points_cents int64) float64 {
	return float64(points_cents) / 100
}
