package settlement

import (
	"context"
	dao "datamonster/connection"
	"fmt"
)

type Repo struct {
	dao dao.Database
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

func NewRepo(d dao.Database) *Repo {
	return &Repo{dao: d}
}

func (r Repo) GetAllForUser(ctx context.Context, userId int) ([]Settlement, error) {
	query := `SELECT * FROM campaign.settlement WHERE owner = $1`
	rows, err := r.dao.Query(ctx, query, userId)
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

func (r Repo) Get(ctx context.Context, id string) (Settlement, error) {
	query := `SELECT * FROM campaign.settlement WHERE id = $1 LIMIT 1`
	var s Settlement
	err := r.dao.QueryRow(ctx, query, id).Scan(&s.Id, &s.Owner, &s.Name, &s.SurvivalLimit, &s.DepartingSurvival, &s.CollectiveCognition, &s.CurrentYear)
	return s, err
}

func (r Repo) Insert(ctx context.Context, s Settlement) (int, error) {
	insert := "INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) "
	values := fmt.Sprintf("VALUES ('%d', '%s', %d, %d, %d, %d) ", s.Owner, s.Name, s.SurvivalLimit, s.DepartingSurvival, s.CollectiveCognition, s.CurrentYear)
	returning := "RETURNING id"
	query := insert + values + returning
	id := 0
	err := r.dao.QueryRow(ctx, query).Scan(&id)
	return id, err
}
