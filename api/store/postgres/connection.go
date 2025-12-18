package postgres

import (
	"context"
	"fmt"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnectionPool(ctx context.Context, dbsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse db config: %w", err)
	}
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return pool, nil
}
