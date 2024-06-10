package survivor

import (
	"context"
	"fmt"
	"github.com/failuretoload/datamonster/store"
	"log"
)

type postgresRepo struct {
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

func newRepo(d store.Connection) *postgresRepo {
	return &postgresRepo{pool: d}
}

func getAllSurvivorsForSettlement(settlementId int) string {
	return fmt.Sprintf("SELECT * FROM campaign.survivor WHERE settlement = %d", settlementId)
}

func (r postgresRepo) Find(ctx context.Context, query string) ([]Survivor, error) {
	log.Default().Println(query)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Default().Println(err.Error())
		return []Survivor{}, err
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
