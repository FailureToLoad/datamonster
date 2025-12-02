package repo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/failuretoload/datamonster/store"
	"github.com/failuretoload/datamonster/survivor/domain"
	"github.com/jackc/pgx/v5"
)

const (
	Table            = "campaign.survivor"
	ID               = "id"
	Settlement       = "settlement"
	Name             = "name"
	Birth            = "birth"
	Gender           = "gender"
	Status           = "status"
	HuntXp           = "huntxp"
	SurvivalColumn   = "survival"
	Movement         = "movement"
	Accuracy         = "accuracy"
	Strength         = "strength"
	Evasion          = "evasion"
	Luck             = "luck"
	Speed            = "speed"
	Insanity         = "insanity"
	SystemicPressure = "systemic_pressure"
	Torment          = "torment"
	Lumi             = "lumi"
	Courage          = "courage"
	Understanding    = "understanding"
)

type Repo struct {
	conn store.Connection
}

func New(c store.Connection) *Repo {
	return &Repo{conn: c}
}

func (r Repo) CreateSurvivor(ctx context.Context, s domain.Survivor) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)",
		Table, Settlement, Name, Birth, Gender, HuntXp, SurvivalColumn, Movement, Accuracy, Strength, Evasion, Luck, Speed, Insanity, SystemicPressure, Torment, Lumi, Courage, Understanding,
	)

	_, err := r.conn.Exec(ctx, query,
		s.Settlement,
		s.Name,
		s.Birth,
		s.Gender,
		s.HuntXp,
		s.Survival,
		s.Movement,
		s.Accuracy,
		s.Strength,
		s.Evasion,
		s.Luck,
		s.Speed,
		s.Insanity,
		s.SystemicPressure,
		s.Torment,
		s.Lumi,
		s.Courage,
		s.Understanding,
	)
	if err != nil {
		logString := fmt.Errorf("unable to create survivor in settlement %d: %w", s.Settlement, err)
		log.Default().Println(logString)

		if IsDuplicateKeyError(err) {
			return NewDuplicateNameError(s.Name)
		}
		return err
	}
	return nil
}

func (r Repo) GetAllSurvivorsForSettlement(ctx context.Context, settlementID int) ([]domain.Survivor, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", Table, Settlement)

	rows, err := r.conn.Query(ctx, query, settlementID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	survivors, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Survivor])
	if err != nil {
		log.Default().Println(err.Error())
		return nil, err
	}

	return survivors, nil
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
