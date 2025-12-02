package repo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/failuretoload/datamonster/survivor/domain"
	"github.com/failuretoload/datamonster/survivor/internal/repo"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurvivorRepo(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := repo.New(mock)
	ctx := context.Background()

	t.Run("create and get survivor", func(t *testing.T) {
		survivor := domain.Survivor{
			Settlement:       1,
			Name:             "Test Survivor",
			Birth:            1,
			Gender:           "F",
			HuntXp:           0,
			Survival:         1,
			Movement:         5,
			Accuracy:         5,
			Strength:         5,
			Evasion:          5,
			Luck:             5,
			Speed:            5,
			Insanity:         0,
			SystemicPressure: 0,
			Torment:          0,
			Lumi:             0,
			Courage:          0,
			Understanding:    0,
		}

		mock.ExpectExec("INSERT INTO campaign.survivor").
			WithArgs(
				1, "Test Survivor", 1, "F", 0, 1, 5, 5, 5, 5, 5, 5, 0, 0, 0, 0, 0, 0,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := db.CreateSurvivor(ctx, survivor)
		require.NoError(t, err)
	})

	t.Run("duplicate name error", func(t *testing.T) {
		survivor := domain.Survivor{
			Settlement: 1,
			Name:       "Unique Name",
			Gender:     "M",
		}

		mock.ExpectExec("INSERT INTO campaign.survivor").
			WithArgs(
				1, "Unique Name", 0, "M", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			).
			WillReturnError(fmt.Errorf("UNIQUE constraint failed"))

		err := db.CreateSurvivor(ctx, survivor)
		assert.Error(t, err)
		assert.IsType(t, repo.DuplicateNameError{}, err)
	})

	t.Run("get survivors from empty settlement", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM campaign.survivor WHERE").
			WithArgs(999).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "settlement", "name", "birth", "gender",
				"status", "huntxp", "survival", "movement",
				"accuracy", "strength", "evasion", "luck",
				"speed", "insanity", "systemic_pressure",
				"torment", "lumi", "courage", "understanding",
			}))

		survivors, err := db.GetAllSurvivorsForSettlement(ctx, 999)
		require.NoError(t, err)
		assert.Empty(t, survivors)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
