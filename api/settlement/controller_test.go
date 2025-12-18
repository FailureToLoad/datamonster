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
	requester.Authorized()
	requester.ExpectUserID("user-with-no-settlements")

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestGetSettlements_ReturnsUserSettlements(t *testing.T) {
	userID := "test-user-1"
	ctx := context.Background()

	_, err := dbContainer.PGPool.Exec(ctx, `
		INSERT INTO settlement (owner, name, survival_limit, departing_survival, collective_cognition, year)
		VALUES ($1, 'First Settlement', 1, 1, 0, 1),
		       ($1, 'Second Settlement', 2, 2, 1, 5)
	`, userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		dbContainer.PGPool.Exec(ctx, "DELETE FROM settlement WHERE owner = $1", userID)
	})

	requester.Authorized()
	requester.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var settlements []domain.Settlement
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settlements))
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
	ctx := context.Background()
	user1 := "isolation-user-1"
	user2 := "isolation-user-2"

	_, err := dbContainer.PGPool.Exec(ctx, `
		INSERT INTO settlement (owner, name, survival_limit, departing_survival, collective_cognition, year)
		VALUES ($1, 'User1 Settlement', 1, 1, 0, 1),
		       ($2, 'User2 Settlement', 1, 1, 0, 1)
	`, user1, user2)
	require.NoError(t, err)
	t.Cleanup(func() {
		dbContainer.PGPool.Exec(ctx, "DELETE FROM settlement WHERE owner IN ($1, $2)", user1, user2)
	})

	requester.Authorized()
	requester.ExpectUserID(user1)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var settlements []domain.Settlement
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settlements))
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
	ctx := context.Background()
	t.Cleanup(func() {
		dbContainer.PGPool.Exec(ctx, "DELETE FROM settlement WHERE owner = $1", userID)
	})

	requester.Authorized()
	requester.ExpectUserID(userID)

	body := `{"name":"New Settlement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var settlementID uuid.UUID
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settlementID))
	assert.NotEqual(t, uuid.Nil, settlementID)

	var name string
	err := dbContainer.PGPool.QueryRow(ctx, "SELECT name FROM settlement WHERE external_id = $1", settlementID).Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "New Settlement", name)
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
	ctx := context.Background()

	var externalID uuid.UUID
	err := dbContainer.PGPool.QueryRow(ctx, `
		INSERT INTO settlement (owner, name, survival_limit, departing_survival, collective_cognition, year)
		VALUES ($1, 'Test Settlement', 5, 3, 2, 10)
		RETURNING external_id
	`, userID).Scan(&externalID)
	require.NoError(t, err)
	t.Cleanup(func() {
		dbContainer.PGPool.Exec(ctx, "DELETE FROM settlement WHERE owner = $1", userID)
	})

	requester.Authorized()
	requester.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+externalID.String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var s domain.Settlement
	require.NoError(t, json.NewDecoder(w.Body).Decode(&s))
	assert.Equal(t, externalID, s.ID)
	assert.Equal(t, "Test Settlement", s.Name)
	assert.Equal(t, 5, s.SurvivalLimit)
	assert.Equal(t, 3, s.DepartingSurvival)
	assert.Equal(t, 2, s.CollectiveCognition)
	assert.Equal(t, 10, s.CurrentYear)
}

func TestGetSettlement_NotFound(t *testing.T) {
	userID := "get-notfound-user"

	requester.Authorized()
	requester.ExpectUserID(userID)

	nonExistentID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+nonExistentID.String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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
	ctx := context.Background()
	user1 := "get-isolation-user-1"
	user2 := "get-isolation-user-2"

	var user2SettlementID uuid.UUID
	err := dbContainer.PGPool.QueryRow(ctx, `
		INSERT INTO settlement (owner, name, survival_limit, departing_survival, collective_cognition, year)
		VALUES ($1, 'User2 Private Settlement', 1, 1, 0, 1)
		RETURNING external_id
	`, user2).Scan(&user2SettlementID)
	require.NoError(t, err)
	t.Cleanup(func() {
		dbContainer.PGPool.Exec(ctx, "DELETE FROM settlement WHERE owner IN ($1, $2)", user1, user2)
	})

	requester.Authorized()
	requester.ExpectUserID(user1)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+user2SettlementID.String(), nil)
	w := httptest.NewRecorder()

	requester.DoRequest(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
