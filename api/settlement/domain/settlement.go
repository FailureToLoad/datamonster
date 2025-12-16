package domain

import "github.com/gofrs/uuid/v5"

type Settlement struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Owner               string    `json:"owner"`
	SurvivalLimit       int       `json:"survivalLimit"`
	DepartingSurvival   int       `json:"departingSurvival"`
	CollectiveCognition int       `json:"collectiveCognition"`
	CurrentYear         int       `json:"currentYear"`
}
