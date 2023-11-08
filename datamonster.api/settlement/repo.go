package settlement

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

type Settlement struct {
	Id                  uuid.UUID
	Owner               uuid.UUID
	Name                string
	SurvivalLimit       int
	DepartingSurvival   int
	CollectiveCognition int
	CurrentYear         int
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}

func (r Repo) GetAllForUser(ctx context.Context, userId string) ([]Settlement, error) {
	query := `SELECT * FROM settlement WHERE owner = $1`
	rows, err := r.conn.Query(ctx, query, userId)
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

func (r Repo) Insert(ctx context.Context, s Settlement) error {
	cmd := `INSERT INTO settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) VALUES (@owner, @name, @limit, @departing, @cc, @year)`
	args := pgx.NamedArgs{
		"owner":     s.Owner,
		"name":      s.Name,
		"limit":     s.SurvivalLimit,
		"departing": s.DepartingSurvival,
		"cc":        s.CollectiveCognition,
		"year":      s.CurrentYear,
	}
	_, err := r.conn.Exec(ctx, cmd, args)
	return err
}
