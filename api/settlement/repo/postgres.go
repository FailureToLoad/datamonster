package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type settlement struct {
	ID                  int    `db:"id"`
	Owner               string `db:"owner"`
	Name                string `db:"name"`
	SurvivalLimit       int    `db:"survival_limit"`
	DepartingSurvival   int    `db:"departing_survival"`
	CollectiveCognition int    `db:"collective_cognition"`
	CurrentYear         int    `db:"year"`
}

const (
	table = "campaign.settlement"
	owner = "owner"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) (*Postgres, error) {
	if p == nil {
		return nil, errors.New("settlement repo: pgx connection pool is required")
	}
	return &Postgres{db: p}, nil
}

func (r Postgres) All(ctx context.Context, userID string) ([]domain.Settlement, error) {
	query := fmt.Sprintf("SELECT * FROM %s where %s = $1", table, owner)

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settlements, err := pgx.CollectRows(rows, pgx.RowToStructByName[settlement])
	if err != nil {
		return nil, err
	}

	if len(settlements) == 0 {
		return nil, nil
	}

	return convertList(settlements), nil
}

func convertList(settlements []settlement) []domain.Settlement {
	var settlementDTOs []domain.Settlement
	for _, s := range settlements {
		dto := domain.Settlement{
			ID:                  s.ID,
			Name:                s.Name,
			SurvivalLimit:       s.SurvivalLimit,
			DepartingSurvival:   s.DepartingSurvival,
			CollectiveCognition: s.CollectiveCognition,
			CurrentYear:         s.CurrentYear,
		}
		settlementDTOs = append(settlementDTOs, dto)
	}
	return settlementDTOs
}
