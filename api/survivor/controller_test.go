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

func createTestSettlement(t *testing.T, userID string) uuid.UUID {
	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"name":"Test Settlement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var settlementID uuid.UUID
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settlementID))
	return settlementID
}

func TestGetSurvivors_Empty(t *testing.T) {
	userID := "survivor-empty-user"

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID.String()+"/survivors", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())
}

func createTestSurvivor(t *testing.T, userID string, settlementID uuid.UUID, name string) domain.Survivor {
	requester.Authorized()
	requester.ExpectUserID(userID)

	body := fmt.Sprintf(`{"name":"%s","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`, name)
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(w.Body).Decode(&survivor))
	return survivor
}

func TestGetSurvivors_ReturnsSurvivors(t *testing.T) {
	userID := "survivor-list-user"

	settlementID := createTestSettlement(t, userID)
	createTestSurvivor(t, userID, settlementID, "Survivor One")
	createTestSurvivor(t, userID, settlementID, "Survivor Two")

	requester.Authorized()
	requester.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID.String()+"/survivors", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var survivors []domain.Survivor
	require.NoError(t, json.NewDecoder(w.Body).Decode(&survivors))
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
	requester.Authorized()
	requester.ExpectUserID("invalid-id-user")

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/not-a-uuid/survivors", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"name":"New Survivor","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(w.Body).Decode(&survivor))
	assert.Equal(t, "New Survivor", survivor.Name)
	assert.NotEqual(t, uuid.Nil, survivor.ID)
	assert.Equal(t, settlementID, survivor.Settlement)
}

func TestCreateSurvivor_MissingName(t *testing.T) {
	userID := "missing-name-survivor-user"

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"birth":1,"gender":"M"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSurvivor_InvalidJSON(t *testing.T) {
	userID := "invalid-json-survivor-user"

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateSurvivor_InvalidSettlementID(t *testing.T) {
	requester.Authorized()
	requester.ExpectUserID("invalid-settlement-user")

	body := `{"name":"Test Survivor"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/not-a-uuid/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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

	settlementID := createTestSettlement(t, userID)
	createTestSurvivor(t, userID, settlementID, "Duplicate Name")

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"name":"Duplicate Name","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpsertSurvivor_CreateNew(t *testing.T) {
	userID := "upsert-create-user"

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"name":"Upsert New Survivor","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(w.Body).Decode(&survivor))
	assert.Equal(t, "Upsert New Survivor", survivor.Name)
	assert.NotEqual(t, uuid.Nil, survivor.ID)
	assert.Equal(t, settlementID, survivor.Settlement)
}

func TestUpsertSurvivor_UpdateExisting(t *testing.T) {
	userID := "upsert-update-user"

	settlementID := createTestSettlement(t, userID)
	existing := createTestSurvivor(t, userID, settlementID, "Original Name")

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := fmt.Sprintf(`{"id":"%s","name":"Original Name","birth":5,"gender":"M","huntxp":10,"survival":8,"movement":7,"accuracy":3,"strength":4,"evasion":2,"luck":5,"speed":3,"insanity":6,"systemicPressure":2,"torment":3,"lumi":4,"courage":7,"understanding":9}`, existing.ID.String())
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var survivor domain.Survivor
	require.NoError(t, json.NewDecoder(w.Body).Decode(&survivor))
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

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"birth":1,"gender":"M"}`
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpsertSurvivor_InvalidJSON(t *testing.T) {
	userID := "upsert-invalid-json-user"

	settlementID := createTestSettlement(t, userID)

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpsertSurvivor_InvalidSettlementID(t *testing.T) {
	requester.Authorized()
	requester.ExpectUserID("upsert-invalid-settlement-user")

	body := `{"name":"Test Survivor"}`
	req := httptest.NewRequest(http.MethodPut, "/api/settlements/not-a-uuid/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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
