package domain

import "github.com/gofrs/uuid/v5"

type Survivor struct {
	ID               uuid.UUID `json:"id"`
	Settlement       uuid.UUID `json:"settlement"`
	Name             string    `json:"name"`
	Birth            int       `json:"birth"`
	Gender           string    `json:"gender"`
	HuntXP           int       `json:"huntxp"`
	Survival         int       `json:"survival"`
	Movement         int       `json:"movement"`
	Accuracy         int       `json:"accuracy"`
	Strength         int       `json:"strength"`
	Evasion          int       `json:"evasion"`
	Luck             int       `json:"luck"`
	Speed            int       `json:"speed"`
	Insanity         int       `json:"insanity"`
	SystemicPressure int       `json:"systemicPressure"`
	Torment          int       `json:"torment"`
	Lumi             int       `json:"lumi"`
	Courage          int       `json:"courage"`
	Understanding    int       `json:"understanding"`
}
