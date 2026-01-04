package migrator

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migration func(context.Context, pgx.Tx) error

var migrations = map[int]migration{
	1: createSettlementTable,
	2: createSurvivorTable,
	3: addFightingArtsToSurvivor,
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = ensureMigrationTable(ctx, pool)
	if err != nil {
		return err
	}

	var lastApplied int
	_ = tx.QueryRow(ctx, "SELECT COALESCE(MAX(id), 0) FROM migration").Scan(&lastApplied)

	for id := 1; id <= len(migrations); id++ {
		if id <= lastApplied {
			continue
		}

		apply := migrations[id]

		if err := apply(ctx, tx); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				slog.Error("did not roll back migration: %w", slog.Any("error", txErr))
			}

			return fmt.Errorf("did not apply migration %d: %w", id, err)
		}

		if recordErr := record(ctx, tx, id); recordErr != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				slog.Error("did not roll back migration: %w", slog.Any("error", txErr))
			}

			return fmt.Errorf("did not record migration %d: %w", id, recordErr)
		}
	}

	return tx.Commit(ctx)
}

func ensureMigrationTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS migration (id INT PRIMARY KEY, applied TIMESTAMPTZ NOT NULL DEFAULT NOW());`)
	if err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	return nil
}

func record(ctx context.Context, tx pgx.Tx, id int) error {
	_, err := tx.Exec(ctx, "INSERT INTO migration (id, applied) VALUES ($1, NOW())", id)
	return err
}
