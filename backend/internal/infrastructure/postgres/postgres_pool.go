package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New_postgres_pool(application_context context.Context, postgres_url string) (*pgxpool.Pool, error) {
	postgres_pool, create_pool_error := pgxpool.New(application_context, postgres_url)
	if create_pool_error != nil {
		return nil, fmt.Errorf("failed_to_create_postgres_pool: %w", create_pool_error)
	}

	ping_database_error := postgres_pool.Ping(application_context)
	if ping_database_error != nil {
		postgres_pool.Close()
		return nil, fmt.Errorf("failed_to_ping_postgres: %w", ping_database_error)
	}

	return postgres_pool, nil
}
