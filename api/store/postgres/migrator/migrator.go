package migrator

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/failuretoload/datamonster/store/postgres/migrator/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migration interface {
	ID() int
	Apply(ctx context.Context, tx pgx.Tx) error
}

var migrations = []migration{
	internal.CreateSettlementTable{},
	internal.CreateSurvivorTable{},
	internal.AddSurvivorStatus{},
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

	for _, m := range migrations {
		mID := m.ID()
		if mID <= lastApplied {
			continue
		}

		if err := m.Apply(ctx, tx); err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				slog.Error("did not roll back migration: %w", slog.Any("error", txErr))
			}

			return fmt.Errorf("did not apply migration %d: %w", mID, err)
		}

		if recordErr := record(ctx, tx, m.ID()); recordErr != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				slog.Error("did not roll back migration: %w", slog.Any("error", txErr))
			}

			return fmt.Errorf("did not record migration %d: %w", mID, recordErr)
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
