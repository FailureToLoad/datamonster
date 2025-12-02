package domain

type Settlement struct {
	ID                  int    `db:"id"`
	Owner               string `db:"owner"`
	Name                string `db:"name"`
	SurvivalLimit       int    `db:"survival_limit"`
	DepartingSurvival   int    `db:"departing_survival"`
	CollectiveCognition int    `db:"collective_cognition"`
	CurrentYear         int    `db:"year"`
}
