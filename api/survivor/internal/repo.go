package internal

import (
	"context"
	"fmt"
	"github.com/failuretoload/datamonster/store"
	"log"
)

type PostGresRepo struct {
	pool store.Connection
}

type Survivor struct {
	Id               int
	Settlement       int
	Name             string
	Born             int
	Gender           string
	Status           string
	HuntXp           int
	Survival         int
	Movement         int
	Accuracy         int
	Strength         int
	Evasion          int
	Luck             int
	Speed            int
	Insanity         int
	SystemicPressure int
	Torment          int
	Lumi             int
	Courage          int
	Understanding    int
}

func NewRepo(d store.Connection) *PostGresRepo {
	return &PostGresRepo{pool: d}
}

func (r PostGresRepo) GetAllSurvivorsForSettlement(ctx context.Context, settlementId int) ([]Survivor, error) {
	query := fmt.Sprintf("SELECT * FROM campaign.survivor WHERE settlement = %d", settlementId)
	survivors, err := r.find(ctx, query)
	return survivors, err
}

func (r PostGresRepo) find(ctx context.Context, query string) ([]Survivor, error) {
	log.Default().Println(query)
	rows, queryErr := r.pool.Query(ctx, query)
	if queryErr != nil {
		log.Default().Println(queryErr.Error())
		return nil, queryErr
	}
	defer rows.Close()
	survivors := []Survivor{}
	for rows.Next() {
		var s Survivor
		err := rows.Scan(&s.Id,
			&s.Settlement,
			&s.Name,
			&s.Born,
			&s.Gender,
			&s.Status,
			&s.HuntXp,
			&s.Survival,
			&s.Movement,
			&s.Accuracy,
			&s.Strength,
			&s.Evasion,
			&s.Luck,
			&s.Speed,
			&s.Insanity,
			&s.SystemicPressure,
			&s.Torment,
			&s.Lumi,
			&s.Courage,
			&s.Understanding,
		)
		if err != nil {
			log.Default().Println(err.Error())
			return survivors, err
		}
		survivors = append(survivors, s)
	}
	return survivors, nil
}
