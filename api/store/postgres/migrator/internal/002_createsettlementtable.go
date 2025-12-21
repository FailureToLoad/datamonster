package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CreateSettlementTable struct{}

func (CreateSettlementTable) ID() int { return 2 }

func (CreateSettlementTable) Apply(ctx context.Context, tx pgx.Tx) error {
	create := `
		CREATE TABLE IF NOT EXISTS settlement (
			id SERIAL PRIMARY KEY,
			external_id UUID NOT NULL DEFAULT gen_random_uuid(),
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
