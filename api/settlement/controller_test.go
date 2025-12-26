package settlement_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/settlement"
	"github.com/failuretoload/datamonster/settlement/domain"
	"github.com/failuretoload/datamonster/settlement/repo"
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
		log.Fatalf("unable to set up test env for settlement tests: %v", err)
	}
	defer dbContainer.Cleanup()

	repo, err := repo.New(dbContainer.PGPool)
	if err != nil {
		log.Fatal(err)
	}

	controller, err := settlement.NewController(repo)
	if err != nil {
		log.Fatal(err)
	}

	requester, err = testenv.NewRequester([]server.Controller{controller})
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetSettlements_Empty(t *testing.T) {
	userID := "user-with-no-settlements"
	body := requester.GetSettlements(t, userID)
	assert.Equal(t, "[]", body.String())
}

func TestGetSettlements_ReturnsUserSettlements(t *testing.T) {
	userID := "test-user-1"

	requester.CreateSettlementWithName(t, userID, "First Settlement")
	requester.CreateSettlementWithName(t, userID, "Second Settlement")
	body := requester.GetSettlements(t, userID)

	var settlements []domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&settlements))
	require.Len(t, settlements, 2)

	names := map[string]bool{}
	for _, s := range settlements {
		names[s.Name] = true
		assert.NotEqual(t, uuid.Nil, s.ID, "expected settlement to have a valid external_id")
	}
	assert.True(t, names["First Settlement"], "expected 'First Settlement' in results")
	assert.True(t, names["Second Settlement"], "expected 'Second Settlement' in results")
}

func TestGetSettlements_IsolatesUserData(t *testing.T) {
	user1 := "isolation-user-1"
	user2 := "isolation-user-2"

	requester.CreateSettlementWithName(t, user1, "User1 Settlement")
	requester.CreateSettlementWithName(t, user2, "User2 Settlement")
	body := requester.GetSettlements(t, user1)

	var settlements []domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&settlements))
	require.Len(t, settlements, 1)
	assert.Equal(t, "User1 Settlement", settlements[0].Name)
	assert.NotEqual(t, uuid.Nil, settlements[0].ID, "expected settlement to have a valid external_id")
}

func TestGetSettlements_Unauthorized(t *testing.T) {
	requester.Unauthorized()

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateSettlement_Success(t *testing.T) {
	userID := "create-test-user"

	settlementID := requester.CreateSettlementWithName(t, userID, "New Settlement")
	getSettlementBody := requester.GetSettlement(t, userID, settlementID)

	var settlement domain.Settlement
	require.NoError(t, json.NewDecoder(getSettlementBody).Decode(&settlement))
	assert.Equal(t, "New Settlement", settlement.Name)
}

func TestCreateSettlement_MissingName(t *testing.T) {
	requester.Authorized()
	requester.ExpectUserID("missing-name-user")

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSettlement_InvalidJSON(t *testing.T) {
	requester.Authorized()
	requester.ExpectUserID("invalid-json-user")

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSettlement_Unauthorized(t *testing.T) {
	requester.Unauthorized()

	body := `{"name":"Unauthorized Settlement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetSettlement_Success(t *testing.T) {
	userID := "get-single-user"

	settlementID := requester.CreateSettlementWithName(t, userID, "Test Settlement")
	body := requester.GetSettlement(t, userID, settlementID)

	var s domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&s))
	assert.Equal(t, settlementID, s.ID)
	assert.Equal(t, "Test Settlement", s.Name)
}

func TestGetSettlement_NotFound(t *testing.T) {
	userID := "get-notfound-user"

	requester.Authorized()
	requester.ExpectUserID(userID)

	nonExistentID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+nonExistentID.String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSettlement_InvalidID(t *testing.T) {
	requester.Authorized()
	requester.ExpectUserID("invalid-id-user")

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/not-a-uuid", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSettlement_Unauthorized(t *testing.T) {
	requester.Unauthorized()

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+uuid.Must(uuid.NewV4()).String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetSettlement_IsolatesUserData(t *testing.T) {
	user1 := "get-isolation-user-1"
	user2 := "get-isolation-user-2"

	user2SettlementID := requester.CreateSettlementWithName(t, user2, "User2 Private Settlement")

	requester.Authorized()
	requester.ExpectUserID(user1)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+user2SettlementID.String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
