package migrator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func createSettlementTable(ctx context.Context, tx pgx.Tx) error {
	create := `
		CREATE TABLE IF NOT EXISTS settlement (
			id SERIAL PRIMARY KEY,
			external_id UUID NOT NULL UNIQUE DEFAULT uuidv7(),
			owner VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			survival_limit INTEGER NOT NULL,
			departing_survival INTEGER NOT NULL,
			collective_cognition INTEGER NOT NULL,
			year INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_settlement_owner ON settlement(owner);
	`

	_, err := tx.Exec(ctx, create)
	if err != nil {
		return fmt.Errorf("failed to create settlement table: %w", err)
	}

	return nil
}

func createSurvivorTable(ctx context.Context, tx pgx.Tx) error {
	create := `
		CREATE TYPE survivor_status AS ENUM ('Alive', 'Ceased to exist', 'Cannot depart', 'Dead', 'Retired');

		CREATE TABLE IF NOT EXISTS survivor (
			id SERIAL PRIMARY KEY,
			external_id UUID NOT NULL DEFAULT uuidv7(),
			settlement_id UUID REFERENCES settlement(external_id),
			name VARCHAR(255) NOT NULL,
			birth INTEGER NOT NULL,
			gender VARCHAR(50) NOT NULL,
			hunt_xp INTEGER NOT NULL,
			survival INTEGER NOT NULL,
			movement INTEGER NOT NULL,
			accuracy INTEGER NOT NULL,
			strength INTEGER NOT NULL,
			evasion INTEGER NOT NULL,
			luck INTEGER NOT NULL,
			speed INTEGER NOT NULL,
			insanity INTEGER NOT NULL,
			systemic_pressure INTEGER NOT NULL,
			torment INTEGER NOT NULL,
			lumi INTEGER NOT NULL,
			courage INTEGER NOT NULL,
			understanding INTEGER NOT NULL,
			status survivor_status NOT NULL DEFAULT 'Alive',
			disorders UUID[3]
		);

		CREATE INDEX IF NOT EXISTS idx_survivors_settlement ON survivor(settlement_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_survivors_settlement_name ON survivor(settlement_id, name);
	`

	_, err := tx.Exec(ctx, create)
	if err != nil {
		return fmt.Errorf("failed to create survivor table: %w", err)
	}

	return nil
}

func addFightingArtsToSurvivor(ctx context.Context, tx pgx.Tx) error {
	alter := `
		ALTER TABLE survivor ADD COLUMN IF NOT EXISTS fighting_art UUID;
		ALTER TABLE survivor ADD COLUMN IF NOT EXISTS secret_fighting_art UUID;
	`

	_, err := tx.Exec(ctx, alter)
	if err != nil {
		return fmt.Errorf("failed to add fighting arts columns to survivor table: %w", err)
	}

	return nil
}
