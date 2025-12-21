package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/failuretoload/datamonster/logger"
	"github.com/failuretoload/datamonster/survivor/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	table            = "survivor"
	settlementID     = "settlement_id"
	name             = "name"
	birth            = "birth"
	gender           = "gender"
	huntXP           = "hunt_xp"
	survival         = "survival"
	movement         = "movement"
	accuracy         = "accuracy"
	strength         = "strength"
	evasion          = "evasion"
	luck             = "luck"
	speed            = "speed"
	insanity         = "insanity"
	systemicPressure = "systemic_pressure"
	torment          = "torment"
	lumi             = "lumi"
	courage          = "courage"
	understanding    = "understanding"
)

func ErrDuplicateName(name string) error {
	return fmt.Errorf("survivor with name %s already exists", name)
}

type Postgres struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) (*Postgres, error) {
	if p == nil {
		return nil, errors.New("survivor repo: pgx connection pool is required")
	}
	return &Postgres{db: p}, nil
}

func (r Postgres) Create(ctx context.Context, d domain.Survivor) (domain.Survivor, error) {
	s := fromDTO(d)
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING *`,
		table,
		settlementID,
		name,
		birth,
		gender,
		huntXP,
		survival,
		movement,
		accuracy,
		strength,
		evasion,
		luck,
		speed,
		insanity,
		systemicPressure,
		torment,
		lumi,
		courage,
		understanding,
	)

	rows, err := r.db.Query(ctx, query,
		s.Settlement,
		s.Name,
		s.Birth,
		s.Gender,
		s.HuntXP,
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
		if IsDuplicateKeyError(err) {
			logger.Error(ctx, fmt.Sprintf("survivor named %s exists", s.Name),
				logger.SettlementID(s.Settlement.String()),
			)
			return domain.Survivor{}, ErrDuplicateName(s.Name)
		}

		safeErr := fmt.Errorf("unable to create survivor")
		logger.Error(ctx, safeErr.Error(),
			logger.SettlementID(s.Settlement.String()),
			logger.ErrorField(err),
		)
		return domain.Survivor{}, safeErr
	}

	inserted, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[survivor])
	if err != nil {
		safeErr := fmt.Errorf("unable to read creation result")
		logger.Error(ctx, safeErr.Error(),
			logger.SettlementID(s.Settlement.String()),
			logger.ErrorField(err),
		)

		return domain.Survivor{}, safeErr
	}

	return toDTO(inserted), nil
}

func (r Postgres) All(ctx context.Context, settlement uuid.UUID) ([]domain.Survivor, error) {
	query := fmt.Sprintf("SELECT * FROM %s where %s = $1", table, settlementID)

	rows, err := r.db.Query(ctx, query, settlement)
	if err != nil {
		safeErr := fmt.Errorf("unable to query survivors for settlement")
		logger.Error(ctx, safeErr.Error(),
			logger.SettlementID(settlement.String()),
			logger.ErrorField(err),
		)
		return nil, safeErr
	}
	defer rows.Close()

	survivors, err := pgx.CollectRows(rows, pgx.RowToStructByName[survivor])
	if err != nil {
		safeErr := fmt.Errorf("unable to scan survivors for settlement")
		logger.Error(ctx, safeErr.Error(),
			logger.SettlementID(settlement.String()),
			logger.ErrorField(err),
		)
		return nil, safeErr
	}

	if len(survivors) == 0 {
		return nil, nil
	}

	return toDTOList(survivors), nil
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

type survivor struct {
	ID               int       `db:"id"`
	ExternalID       uuid.UUID `db:"external_id"`
	Settlement       uuid.UUID `db:"settlement_id"`
	Name             string    `db:"name"`
	Birth            int       `db:"birth"`
	Gender           string    `db:"gender"`
	HuntXP           int       `db:"hunt_xp"`
	Survival         int       `db:"survival"`
	Movement         int       `db:"movement"`
	Accuracy         int       `db:"accuracy"`
	Strength         int       `db:"strength"`
	Evasion          int       `db:"evasion"`
	Luck             int       `db:"luck"`
	Speed            int       `db:"speed"`
	Insanity         int       `db:"insanity"`
	SystemicPressure int       `db:"systemic_pressure"`
	Torment          int       `db:"torment"`
	Lumi             int       `db:"lumi"`
	Courage          int       `db:"courage"`
	Understanding    int       `db:"understanding"`
}

func toDTO(s survivor) domain.Survivor {
	return domain.Survivor{
		ID:               s.ExternalID,
		Settlement:       s.Settlement,
		Name:             s.Name,
		Birth:            s.Birth,
		Gender:           s.Gender,
		HuntXP:           s.HuntXP,
		Survival:         s.Survival,
		Movement:         s.Movement,
		Accuracy:         s.Accuracy,
		Strength:         s.Strength,
		Evasion:          s.Evasion,
		Luck:             s.Luck,
		Speed:            s.Speed,
		Insanity:         s.Insanity,
		SystemicPressure: s.SystemicPressure,
		Torment:          s.Torment,
		Lumi:             s.Lumi,
		Courage:          s.Courage,
		Understanding:    s.Understanding,
	}
}

func toDTOList(survivors []survivor) []domain.Survivor {
	dtos := make([]domain.Survivor, len(survivors))

	for i, s := range survivors {
		dtos[i] = toDTO(s)
	}

	return dtos
}

func fromDTO(s domain.Survivor) survivor {
	return survivor{
		ExternalID:       s.ID,
		Settlement:       s.Settlement,
		Name:             s.Name,
		Birth:            s.Birth,
		Gender:           s.Gender,
		HuntXP:           s.HuntXP,
		Survival:         s.Survival,
		Movement:         s.Movement,
		Accuracy:         s.Accuracy,
		Strength:         s.Strength,
		Evasion:          s.Evasion,
		Luck:             s.Luck,
		Speed:            s.Speed,
		Insanity:         s.Insanity,
		SystemicPressure: s.SystemicPressure,
		Torment:          s.Torment,
		Lumi:             s.Lumi,
		Courage:          s.Courage,
		Understanding:    s.Understanding,
	}
}
