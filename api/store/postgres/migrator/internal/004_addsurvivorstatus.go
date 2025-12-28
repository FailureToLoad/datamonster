package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type AddSurvivorStatus struct{}

func (AddSurvivorStatus) ID() int { return 4 }

func (AddSurvivorStatus) Apply(ctx context.Context, tx pgx.Tx) error {
	alter := `ALTER TABLE survivor ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'Alive';`

	_, err := tx.Exec(ctx, alter)
	if err != nil {
		return fmt.Errorf("failed to add status column to survivor table: %w", err)
	}

	return nil
}
