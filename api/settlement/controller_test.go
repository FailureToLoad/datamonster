package settlement_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	body, status := requester.GetSettlements(userID)
	require.Equal(t, http.StatusOK, status)

	assert.Equal(t, "[]", body.String())
}

func TestGetSettlements_ReturnsUserSettlements(t *testing.T) {
	userID := "test-user-1"

	_, status := requester.CreateSettlementWithBody(userID, `{"name":"First Settlement"}`)
	require.Equal(t, http.StatusOK, status)
	_, status = requester.CreateSettlementWithBody(userID, `{"name":"Second Settlement"}`)
	require.Equal(t, http.StatusOK, status)

	body, status := requester.GetSettlements(userID)
	require.Equal(t, http.StatusOK, status)

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

	_, status := requester.CreateSettlementWithBody(user1, `{"name":"User1 Settlement"}`)
	require.Equal(t, http.StatusOK, status)

	_, status = requester.CreateSettlementWithBody(user2, `{"name":"User2 Settlement"}`)
	require.Equal(t, http.StatusOK, status)

	body, status := requester.GetSettlements(user1)
	require.Equal(t, http.StatusOK, status)

	var settlements []domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&settlements))
	require.Len(t, settlements, 1)
	assert.Equal(t, "User1 Settlement", settlements[0].Name)
	assert.NotEqual(t, uuid.Nil, settlements[0].ID, "expected settlement to have a valid external_id")
}

func TestGetSettlements_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.GetSettlements("unauthorized")

	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestCreateSettlement_Success(t *testing.T) {
	userID := "create-test-user"

	respBody, status := requester.CreateSettlementWithBody(userID, `{"name":"New Settlement"}`)
	require.Equal(t, http.StatusOK, status)

	var settlementID uuid.UUID
	require.NoError(t, json.NewDecoder(respBody).Decode(&settlementID))

	body, status := requester.GetSettlement(userID, settlementID.String())
	require.Equal(t, http.StatusOK, status)

	var settlement domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&settlement))
	assert.Equal(t, "New Settlement", settlement.Name)
}

func TestCreateSettlement_MissingName(t *testing.T) {
	_, status := requester.CreateSettlementWithBody("missing-name-user", `{}`)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestCreateSettlement_InvalidJSON(t *testing.T) {
	_, status := requester.CreateSettlementWithBody("invalid-json-user", `{invalid json}`)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestCreateSettlement_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.CreateSettlementWithBody("unauthorized", `{"name":"Unauthorized Settlement"}`)

	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestGetSettlement_Success(t *testing.T) {
	userID := "get-single-user"

	settlementID, err := requester.CreateSettlement(userID)
	require.NoError(t, err)

	body, status := requester.GetSettlement(userID, settlementID)
	require.Equal(t, http.StatusOK, status)

	var s domain.Settlement
	require.NoError(t, json.NewDecoder(body).Decode(&s))
	assert.Equal(t, settlementID, s.ID.String())
	assert.Equal(t, "Test Settlement", s.Name)
}

func TestGetSettlement_NotFound(t *testing.T) {
	_, status := requester.GetSettlement("get-notfound-user", testenv.ValidUUID())
	assert.Equal(t, http.StatusNotFound, status)
}

func TestGetSettlement_InvalidID(t *testing.T) {
	_, status := requester.GetSettlement("invalid-id-user", "not-a-uuid")
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestGetSettlement_Unauthorized(t *testing.T) {
	t.Cleanup(requester.Unauthorized())
	_, status := requester.GetSettlement("unauthorized", testenv.ValidUUID())

	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestGetSettlement_IsolatesUserData(t *testing.T) {
	user1 := "get-isolation-user-1"
	user2 := "get-isolation-user-2"

	user2SettlementID, err := requester.CreateSettlement(user2)
	require.NoError(t, err)

	_, status := requester.GetSettlement(user1, user2SettlementID)
	assert.Equal(t, http.StatusNotFound, status)
}
