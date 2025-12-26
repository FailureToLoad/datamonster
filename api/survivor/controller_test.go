package survivor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	requester.Unauthorized()

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+uuid.Must(uuid.NewV4()).String()+"/survivors", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
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
	requester.Unauthorized()

	body := `{"name":"Unauthorized Survivor"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+uuid.Must(uuid.NewV4()).String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
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

func TestUpsertSurvivor_CreateNew(t *testing.T) {
	userID := "upsert-create-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	body := `{"name":"Upsert New Survivor","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`
	respBody, status := requester.UpsertSurvivor(userID, settlementID, body)
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, "Upsert New Survivor", survivor.Name)
	assert.NotEqual(t, uuid.Nil, survivor.ID)
}

func TestUpsertSurvivor_UpdateExisting(t *testing.T) {
	userID := "upsert-update-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	rawSurvivor, status := requester.CreateSurvivor(userID, settlementID, "Original Name")
	require.Equal(t, http.StatusOK, status)

	var existing domain.Survivor
	require.NoError(t, json.NewDecoder(rawSurvivor).Decode(&existing))

	body := fmt.Sprintf(`{"id":"%s","name":"Original Name","birth":5,"gender":"M","huntxp":10,"survival":8,"movement":7,"accuracy":3,"strength":4,"evasion":2,"luck":5,"speed":3,"insanity":6,"systemicPressure":2,"torment":3,"lumi":4,"courage":7,"understanding":9}`, existing.ID.String())
	respBody, status := requester.UpsertSurvivor(userID, settlementID, body)
	require.Equal(t, http.StatusOK, status)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(respBody).Decode(&survivor))
	assert.Equal(t, existing.ID, survivor.ID)
	assert.Equal(t, existing.Name, survivor.Name)
	assert.Equal(t, existing.Gender, survivor.Gender)
	assert.Equal(t, 5, survivor.Birth)
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

func TestUpsertSurvivor_MissingName(t *testing.T) {
	userID := "upsert-missing-name-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	_, status := requester.UpsertSurvivor(userID, settlementID, `{"birth":1,"gender":"M"}`)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestUpsertSurvivor_InvalidJSON(t *testing.T) {
	userID := "upsert-invalid-json-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	_, status := requester.UpsertSurvivor(userID, settlementID, `{invalid json}`)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpsertSurvivor_InvalidSettlementID(t *testing.T) {
	_, status := requester.UpsertSurvivor("upsert-invalid-settlement-user", "not-a-uuid", `{"name":"Test Survivor"}`)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestUpsertSurvivor_Unauthorized(t *testing.T) {
	requester.Unauthorized()

	body := `{"name":"Unauthorized Survivor"}`
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+uuid.Must(uuid.NewV4()).String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
