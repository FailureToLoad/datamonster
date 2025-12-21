package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CreateSurvivorTable struct{}

func (CreateSurvivorTable) ID() int { return 3 }

func (CreateSurvivorTable) Apply(ctx context.Context, tx pgx.Tx) error {
	create := `
		ALTER TABLE settlement ADD CONSTRAINT settlement_external_id_unique UNIQUE (external_id);

		CREATE TABLE IF NOT EXISTS survivor (
			id SERIAL PRIMARY KEY,
			external_id UUID NOT NULL DEFAULT gen_random_uuid(),
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
			understanding INTEGER NOT NULL
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
