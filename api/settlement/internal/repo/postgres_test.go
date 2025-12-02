package repo_test

import (
	"context"
	"testing"

	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/failuretoload/datamonster/settlement/internal/repo"
	"github.com/jackc/pgx/v5"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettlementRepo(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := repo.New(mock)
	ctx := context.Background()

	t.Run("insert and get settlement", func(t *testing.T) {
		settlement := domain.Settlement{
			Owner:               "test-user",
			Name:                "Lantern Hoard",
			SurvivalLimit:       3,
			DepartingSurvival:   2,
			CollectiveCognition: 4,
			CurrentYear:         1,
		}

		mock.ExpectQuery("INSERT INTO campaign.settlement").
			WithArgs(settlement.Owner, settlement.Name, settlement.SurvivalLimit,
				settlement.DepartingSurvival, settlement.CollectiveCognition, settlement.CurrentYear).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int32(1)))

		id, err := db.Insert(ctx, settlement)
		require.NoError(t, err)
		assert.Equal(t, 1, id)

		mock.ExpectQuery("SELECT .* FROM campaign.settlement where").
			WithArgs(id, settlement.Owner).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "owner", "name", "survival_limit",
				"departing_survival", "collective_cognition", "year",
			}).AddRow(
				1, settlement.Owner, settlement.Name, settlement.SurvivalLimit,
				settlement.DepartingSurvival, settlement.CollectiveCognition, settlement.CurrentYear,
			))

		retrieved, err := db.Get(ctx, id, settlement.Owner)
		require.NoError(t, err)

		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, settlement.Owner, retrieved.Owner)
		assert.Equal(t, settlement.Name, retrieved.Name)
		assert.Equal(t, settlement.SurvivalLimit, retrieved.SurvivalLimit)
		assert.Equal(t, settlement.DepartingSurvival, retrieved.DepartingSurvival)
		assert.Equal(t, settlement.CollectiveCognition, retrieved.CollectiveCognition)
		assert.Equal(t, settlement.CurrentYear, retrieved.CurrentYear)
	})

	t.Run("get non-existent settlement", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM campaign.settlement where").
			WithArgs(999, "test-user").
			WillReturnError(pgx.ErrNoRows)

		_, err := db.Get(ctx, 999, "test-user")
		assert.Error(t, err)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM campaign.settlement where").
			WithArgs(1, "user-2").
			WillReturnError(pgx.ErrNoRows)

		_, err := db.Get(ctx, 1, "user-2")
		assert.Error(t, err)
	})

	t.Run("all settlements", func(t *testing.T) {
		owner := "test-user-select"

		mock.ExpectQuery("SELECT .* FROM campaign.settlement where").
			WithArgs(owner).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "owner", "name", "survival_limit",
				"departing_survival", "collective_cognition", "year",
			}).
				AddRow(1, owner, "First Settlement", 1, 1, 1, 1).
				AddRow(2, owner, "Second Settlement", 2, 2, 2, 2))

		retrieved, err := db.All(ctx, owner)
		require.NoError(t, err)

		assert.Len(t, retrieved, 2)
		for _, s := range retrieved {
			assert.Equal(t, owner, s.Owner)
		}

		names := make(map[string]bool)
		for _, s := range retrieved {
			names[s.Name] = true
		}
		assert.True(t, names["First Settlement"])
		assert.True(t, names["Second Settlement"])
	})

	t.Run("select returns empty slice when no settlements", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM campaign.settlement where").
			WithArgs("non-existent-user").
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "owner", "name", "survival_limit",
				"departing_survival", "collective_cognition", "year",
			}))

		retrieved, err := db.All(ctx, "non-existent-user")
		require.NoError(t, err)
		assert.Empty(t, retrieved)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
