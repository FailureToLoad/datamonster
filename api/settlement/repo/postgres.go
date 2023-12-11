package repo

import (
	"context"
	"datamonster/store"
	"fmt"
)

type PostgresRepo struct {
	pool store.Connection
}

type Settlement struct {
	Id                  int
	Owner               int
	Name                string
	SurvivalLimit       int
	DepartingSurvival   int
	CollectiveCognition int
	CurrentYear         int
}

func New(d store.Connection) *PostgresRepo {
	return &PostgresRepo{pool: d}
}

func (r PostgresRepo) Select(ctx context.Context, query string) ([]Settlement, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	settlements := []Settlement{}
	for rows.Next() {
		var s Settlement
		err := rows.Scan(&s.Id, &s.Owner, &s.Name, &s.SurvivalLimit, &s.DepartingSurvival, &s.CollectiveCognition, &s.CurrentYear)
		if err != nil {
			return settlements, err
		}
		settlements = append(settlements, s)
	}
	return settlements, nil
}

func (r PostgresRepo) Get(ctx context.Context, id string) (Settlement, error) {
	query := `SELECT * FROM campaign.settlement WHERE id = $1 LIMIT 1`
	var s Settlement
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.Id, &s.Owner, &s.Name, &s.SurvivalLimit, &s.DepartingSurvival, &s.CollectiveCognition, &s.CurrentYear)
	return s, err
}

func (r PostgresRepo) Insert(ctx context.Context, s Settlement) (int, error) {
	insert := "INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) "
	values := fmt.Sprintf("VALUES ('%d', '%s', %d, %d, %d, %d) ", s.Owner, s.Name, s.SurvivalLimit, s.DepartingSurvival, s.CollectiveCognition, s.CurrentYear)
	returning := "RETURNING id"
	query := insert + values + returning
	id := 0
	err := r.pool.QueryRow(ctx, query).Scan(&id)
	return id, err
}
