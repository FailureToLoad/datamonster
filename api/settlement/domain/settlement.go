package domain

type Settlement struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Owner               string `json:"owner"`
	SurvivalLimit       int    `json:"survivalLimit"`
	DepartingSurvival   int    `json:"departingSurvival"`
	CollectiveCognition int    `json:"collectiveCognition"`
	CurrentYear         int    `json:"currentYear"`
}
