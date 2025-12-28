package domain

import "github.com/gofrs/uuid/v5"

type SurvivorStatus string

const (
	StatusAlive         SurvivorStatus = "Alive"
	StatusDead          SurvivorStatus = "Dead"
	StatusCeasedToExist SurvivorStatus = "Ceased to exist"
	StatusCannotDepart  SurvivorStatus = "Cannot depart"
	StatusRetired       SurvivorStatus = "Retired"
)

func ValidStatus(s string) bool {
	switch SurvivorStatus(s) {
	case StatusAlive, StatusDead, StatusCeasedToExist, StatusCannotDepart, StatusRetired:
		return true
	}
	return false
}

type Survivor struct {
	ID               uuid.UUID      `json:"id"`
	SettlementID     uuid.UUID      `json:"settlementId"`
	Name             string         `json:"name"`
	Birth            int            `json:"birth"`
	Gender           string         `json:"gender"`
	Status           SurvivorStatus `json:"status"`
	HuntXP           int            `json:"huntxp"`
	Survival         int            `json:"survival"`
	Movement         int            `json:"movement"`
	Accuracy         int            `json:"accuracy"`
	Strength         int            `json:"strength"`
	Evasion          int            `json:"evasion"`
	Luck             int            `json:"luck"`
	Speed            int            `json:"speed"`
	Insanity         int            `json:"insanity"`
	SystemicPressure int            `json:"systemicPressure"`
	Torment          int            `json:"torment"`
	Lumi             int            `json:"lumi"`
	Courage          int            `json:"courage"`
	Understanding    int            `json:"understanding"`
}

type SurvivorUpdate struct {
	StatUpdates  map[string]int  `json:"statUpdates,omitempty"`
	StatusUpdate *SurvivorStatus `json:"statusUpdate,omitempty"`
}
