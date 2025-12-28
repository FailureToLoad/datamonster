package survivor_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/settlement"
	settlementRepo "github.com/failuretoload/datamonster/settlement/repo"
	"github.com/failuretoload/datamonster/survivor"
	"github.com/failuretoload/datamonster/survivor/domain"
	survivorRepo "github.com/failuretoload/datamonster/survivor/repo"
	"github.com/failuretoload/datamonster/testenv"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dbContainer *testenv.DBContainer
	requester   *testenv.Requester
)

func TestMain(m *testing.M) {
	var err error
	dbContainer, err = testenv.NewDBContainer(context.Background())
	if err != nil {
		log.Fatalf("unable to set up test env for survivor tests: %v", err)
	}
	defer dbContainer.Cleanup()

	settlementRepo, err := settlementRepo.New(dbContainer.PGPool)
	if err != nil {
		log.Fatal(err)
	}
	settlementController, err := settlement.NewController(settlementRepo)
	if err != nil {
		log.Fatal(err)
	}

	survivorRepo, err := survivorRepo.New(dbContainer.PGPool)
	if err != nil {
		log.Fatal(err)
	}
	survivorController, err := survivor.NewController(survivorRepo)
	if err != nil {
		log.Fatal(err)
	}

	requester, err = testenv.NewRequester([]server.Controller{settlementController, survivorController})
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetSurvivors_Empty(t *testing.T) {
	userID := "survivor-empty-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	body, status := requester.GetSurvivors(userID, settlementID)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "null", body.String())
}

func TestGetSurvivors_ReturnsSurvivors(t *testing.T) {
	userID := "survivor-list-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)
	_, survivorOneStatus := requester.CreateSurvivor(userID, settlementID, "Survivor One")
	require.Equal(t, http.StatusOK, survivorOneStatus)

	_, survivorTwoStatus := requester.CreateSurvivor(userID, settlementID, "Survivor Two")
	require.Equal(t, http.StatusOK, survivorTwoStatus)

	body, status := requester.GetSurvivors(userID, settlementID)
	require.Equal(t, http.StatusOK, status)

	var survivors []domain.Survivor
	require.NoError(t, json.NewDecoder(body).Decode(&survivors))
	require.Len(t, survivors, 2)

	names := map[string]bool{}
	for _, s := range survivors {
		names[s.Name] = true
		assert.NotEqual(t, uuid.Nil, s.ID)
	}
	assert.True(t, names["Survivor One"])
	assert.True(t, names["Survivor Two"])
}

func TestGetSurvivors_InvalidSettlementID(t *testing.T) {
	_, status := requester.GetSurvivors("invalid-id-user", "not-a-uuid")
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestGetSurvivors_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.GetSurvivors("unauthorized", testenv.UUIDString())

	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestCreateSurvivor_Success(t *testing.T) {
	userID := "create-survivor-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	respBody, status := requester.CreateSurvivor(userID, settlementID, "New Survivor")
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, "New Survivor", survivor.Name)
	assert.NotEqual(t, uuid.Nil, survivor.ID)
}

func TestCreateSurvivor_MissingName(t *testing.T) {
	userID := "missing-name-survivor-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	_, status := requester.CreateSurvivorWithBody(userID, settlementID, `{"birth":1,"gender":"M"}`)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestCreateSurvivor_InvalidJSON(t *testing.T) {
	userID := "invalid-json-survivor-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	_, status := requester.CreateSurvivorWithBody(userID, settlementID, `{invalid json}`)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestCreateSurvivor_InvalidSettlementID(t *testing.T) {
	_, status := requester.CreateSurvivorWithBody("invalid-settlement-user", "not-a-uuid", `{"name":"Test Survivor"}`)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestCreateSurvivor_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.CreateSurvivorWithBody("unauthorized", testenv.UUIDString(), `{"name":"Unauthorized Survivor"}`)

	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestCreateSurvivor_DuplicateName(t *testing.T) {
	userID := "duplicate-survivor-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	_, status := requester.CreateSurvivor(userID, settlementID, "Duplicate Name")
	require.Equal(t, http.StatusOK, status)

	_, status = requester.CreateSurvivor(userID, settlementID, "Duplicate Name")
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpdateSurvivor_UpdateExisting(t *testing.T) {
	userID := "upsert-update-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	rawSurvivor, status := requester.CreateSurvivor(userID, settlementID, "Original Name")
	require.Equal(t, http.StatusOK, status)

	var existing domain.Survivor
	require.NoError(t, json.NewDecoder(rawSurvivor).Decode(&existing))

	body := `{"statUpdates":{"huntxp":10,"survival":8,"movement":7,"accuracy":3,"strength":4,"evasion":2,"luck":5,"speed":3,"insanity":6,"systemicPressure":2,"torment":3,"lumi":4,"courage":7,"understanding":9}}`
	respBody, status := requester.UpdateSurvivor(userID,
		settlementID,
		existing.ID.String(),
		body,
	)
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, existing.ID, survivor.ID)
	assert.Equal(t, existing.Name, survivor.Name)
	assert.Equal(t, existing.Gender, survivor.Gender)
	assert.Equal(t, existing.Birth, survivor.Birth)
	assert.Equal(t, 10, survivor.HuntXP)
	assert.Equal(t, 8, survivor.Survival)
	assert.Equal(t, 7, survivor.Movement)
	assert.Equal(t, 3, survivor.Accuracy)
	assert.Equal(t, 4, survivor.Strength)
	assert.Equal(t, 2, survivor.Evasion)
	assert.Equal(t, 5, survivor.Luck)
	assert.Equal(t, 3, survivor.Speed)
	assert.Equal(t, 6, survivor.Insanity)
	assert.Equal(t, 2, survivor.SystemicPressure)
	assert.Equal(t, 3, survivor.Torment)
	assert.Equal(t, 4, survivor.Lumi)
	assert.Equal(t, 7, survivor.Courage)
	assert.Equal(t, 9, survivor.Understanding)
}

func TestUpdateSurvivor_UpdateStatus(t *testing.T) {
	userID := "update-status-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	rawSurvivor, status := requester.CreateSurvivor(userID, settlementID, "Status Test")
	require.Equal(t, http.StatusOK, status)

	var existing domain.Survivor
	require.NoError(t, json.NewDecoder(rawSurvivor).Decode(&existing))
	assert.Equal(t, domain.StatusAlive, existing.Status)

	body := `{"statusUpdate":"Dead"}`
	respBody, status := requester.UpdateSurvivor(userID,
		settlementID,
		existing.ID.String(),
		body,
	)
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, domain.StatusDead, survivor.Status)
}

func TestUpdateSurvivor_UpdateStatusAndStats(t *testing.T) {
	userID := "update-status-and-stats-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	rawSurvivor, status := requester.CreateSurvivor(userID, settlementID, "Combined Test")
	require.Equal(t, http.StatusOK, status)

	var existing domain.Survivor
	require.NoError(t, json.NewDecoder(rawSurvivor).Decode(&existing))

	body := `{"statUpdates":{"huntxp":5,"courage":3},"statusUpdate":"Retired"}`
	respBody, status := requester.UpdateSurvivor(userID,
		settlementID,
		existing.ID.String(),
		body,
	)
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, domain.StatusRetired, survivor.Status)
	assert.Equal(t, 5, survivor.HuntXP)
	assert.Equal(t, 3, survivor.Courage)
}

func TestUpdateSurvivor_InvalidStatus(t *testing.T) {
	userID := "update-invalid-status-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	rawSurvivor, status := requester.CreateSurvivor(userID, settlementID, "Invalid Status Test")
	require.Equal(t, http.StatusOK, status)

	var existing domain.Survivor
	require.NoError(t, json.NewDecoder(rawSurvivor).Decode(&existing))

	body := `{"statusUpdate":"NotAValidStatus"}`
	_, status = requester.UpdateSurvivor(userID,
		settlementID,
		existing.ID.String(),
		body,
	)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestUpdateSurvivor_NoValidFields(t *testing.T) {
	userID := "update-no-valid-fields-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)
	survivorID := testenv.UUIDString()

	_, status := requester.UpdateSurvivor(userID,
		settlementID,
		survivorID,
		`{"statUpdates":{"notafield":99}}`,
	)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpdateSurvivor_InvalidSettlementID(t *testing.T) {
	_, status := requester.UpdateSurvivor("upsert-invalid-settlement-user",
		"not-a-uuid",
		testenv.UUIDString(),
		`{"statUpdates":{"huntxp":5}}`,
	)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpdateSurvivor_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.UpdateSurvivor("unauthorized", testenv.UUIDString(), testenv.UUIDString(), `{"statUpdates":{"huntxp":5}}`)

	assert.Equal(t, http.StatusUnauthorized, status)
}
