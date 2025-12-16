package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type settlement struct {
	ID                  int       `db:"id"`
	ExternalID          uuid.UUID `db:"external_id"`
	Owner               string    `db:"owner"`
	Name                string    `db:"name"`
	SurvivalLimit       int       `db:"survival_limit"`
	DepartingSurvival   int       `db:"departing_survival"`
	CollectiveCognition int       `db:"collective_cognition"`
	CurrentYear         int       `db:"year"`
}

const (
	table               = "campaign.settlement"
	owner               = "owner"
	id                  = "id"
	externalID          = "external_id"
	name                = "name"
	survivalLimit       = "survival_limit"
	departingSurvival   = "departing_survival"
	collectiveCognition = "collective_cognition"
	year                = "year"
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

	return toDTOList(settlements), nil
}

func (r Postgres) Insert(ctx context.Context, s domain.Settlement) (int, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6) RETURNING %s",
		table, owner, name, survivalLimit, departingSurvival, collectiveCognition, year, id,
	)

	var id int32
	err := r.db.QueryRow(ctx, query,
		s.Owner,
		s.Name,
		s.SurvivalLimit,
		s.DepartingSurvival,
		s.CollectiveCognition,
		s.CurrentYear,
	).Scan(&id)

	return int(id), err
}

func (r Postgres) Get(ctx context.Context, userID string, settlementID uuid.UUID) (domain.Settlement, error) {
	query := fmt.Sprintf("SELECT * FROM %s where %s = $1 AND %s = $2", table, owner, externalID)

	row, err := r.db.Query(ctx, query, userID, settlementID)
	if err != nil {
		return domain.Settlement{}, err
	}
	defer row.Close()

	s, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[settlement])
	if err != nil {
		return domain.Settlement{}, err
	}

	return toDTO(s), nil
}

func toDTOList(settlements []settlement) []domain.Settlement {
	var settlementDTOs []domain.Settlement
	for _, s := range settlements {
		settlementDTOs = append(settlementDTOs, toDTO(s))
	}
	return settlementDTOs
}

func toDTO(s settlement) domain.Settlement {
	return domain.Settlement{
		ID:                  s.ExternalID,
		Owner:               s.Owner,
		Name:                s.Name,
		SurvivalLimit:       s.SurvivalLimit,
		DepartingSurvival:   s.DepartingSurvival,
		CollectiveCognition: s.CollectiveCognition,
		CurrentYear:         s.CurrentYear,
	}
}
