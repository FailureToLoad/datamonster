package repo

import (
	"context"
	"fmt"

	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/failuretoload/datamonster/store"
	"github.com/jackc/pgx/v5"
)

const (
	Table             = "campaign.settlement"
	ID                = "id"
	Owner             = "owner"
	Name              = "name"
	SurvivalLimit     = "survival_limit"
	DepartingSurvival = "departing_survival"
	CC                = "collective_cognition"
	Year              = "year"
)

type Repo struct {
	conn store.Connection
}

func New(c store.Connection) *Repo {
	return &Repo{conn: c}
}

func (r Repo) All(ctx context.Context, userID string) ([]domain.Settlement, error) {
	query := fmt.Sprintf("SELECT * FROM %s where %s = $1", Table, Owner)

	rows, err := r.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settlements, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Settlement])
	if err != nil {
		return nil, err
	}
	return settlements, nil
}

func (r Repo) Get(ctx context.Context, id int, userID string) (domain.Settlement, error) {
	query := fmt.Sprintf("SELECT * FROM %s where %s = $1 AND %s = $2", Table, ID, Owner)
	rows, err := r.conn.Query(ctx, query, id, userID)
	if err != nil {
		return domain.Settlement{}, err
	}
	defer rows.Close()

	settlement, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Settlement])
	if err != nil {
		return domain.Settlement{}, err
	}
	return settlement, nil
}

func (r Repo) Insert(ctx context.Context, s domain.Settlement) (int, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6) RETURNING %s",
		Table, Owner, Name, SurvivalLimit, DepartingSurvival, CC, Year, ID,
	)

	var id int32
	err := r.conn.QueryRow(ctx, query,
		s.Owner,
		s.Name,
		s.SurvivalLimit,
		s.DepartingSurvival,
		s.CollectiveCognition,
		s.CurrentYear,
	).Scan(&id)
	return int(id), err
}
