package settlement_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/settlement"
	"github.com/pashagolub/pgxmock/v4"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testUserID     = "userId"
	settlementCols = []string{"id", "owner", "name", "survival_limit", "departing_survival", "collective_cognition", "year"}
)

func normalizeSQL(sql string) string {
	// Remove all whitespace between symbols and parentheses
	re := regexp.MustCompile(`\s*([(),])\s*`)
	sql = re.ReplaceAllString(sql, "$1")

	// Collapse multiple spaces into single space
	re = regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(sql, " "))
}

type sqlMatcher struct{}

func (m sqlMatcher) Match(expectedSQL, actualSQL string) error {
	if normalizeSQL(expectedSQL) != normalizeSQL(actualSQL) {
		return fmt.Errorf(`
			SQL does not match
			Expected: %s
			Actual: %s`,
			normalizeSQL(expectedSQL),
			normalizeSQL(actualSQL))
	}
	return nil
}

func setupTest(t *testing.T) (pgxmock.PgxPoolIface, *chi.Mux) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(sqlMatcher{}))
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	controller, err := settlement.NewController(mock)
	require.NoError(t, err)
	router := chi.NewRouter()
	router.Use(request.AutheliaMiddleware)
	controller.RegisterRoutes(router)

	return mock, router
}

func TestGetSettlements_ReturnsSettlementsList(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	values := [][]any{
		{1, testUserID, "Fun Forever", 1, 0, 0, 1},
		{2, testUserID, "Wait, we get insanity for the croc?", 1, 0, 0, 1},
	}
	rows := pgxmock.NewRows(settlementCols).AddRows(values...)
	db.ExpectQuery("SELECT * FROM campaign.settlement where owner = $1").
		WithArgs(testUserID).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/settlements", nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	var dto []settlement.DTO
	_ = json.Unmarshal(body, &dto)
	assert.Equal(t, 2, len(dto), "2 settlements should be returned")
}

func TestGetSettlements_ReportsScanErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement where owner = $1").
		WithArgs(testUserID).
		WillReturnError(fmt.Errorf("scan error"))

	req := httptest.NewRequest("GET", "/settlements", nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "scan errors should result in server error")
}

func TestGetSettlements_ReportsConnectionErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement where owner = $1").
		WithArgs(testUserID).
		WillReturnError(fmt.Errorf("query error"))

	req := httptest.NewRequest("GET", "/settlements", nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "connection issues should result in server error")
}

func TestCreateSettlement_ReturnsASettlement(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id").
		WithArgs(testUserID, "Fun Forever", 1, 0, 0, 1).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int32(1)))

	settlementRequest := settlement.CreateSettlementRequest{
		Name: "Fun Forever",
	}
	reqBody, _ := json.Marshal(settlementRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "return 200 on success")
	respBody, _ := io.ReadAll(resp.Body)
	dto := settlement.DTO{}
	_ = json.Unmarshal(respBody, &dto)
	assert.Equal(t, 1, dto.ID, "created settlement should have an id")
	assert.Equal(t, "Fun Forever", dto.Name, "created settlement should have supplied name")
}

func TestCreateSettlement_EnforceRequestType(t *testing.T) {
	_, router := setupTest(t)

	wrongRequest := WrongRequest{
		FancyName: "Fun Forever",
	}
	reqBody, _ := json.Marshal(wrongRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 400, resp.StatusCode, "Request must be of type CreateSettlementRequest")
}

func TestCreateSettlement_RequiresAName(t *testing.T) {
	_, router := setupTest(t)

	emptyRequest := settlement.CreateSettlementRequest{
		Name: "",
	}
	reqBody, _ := json.Marshal(emptyRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 400, resp.StatusCode, "Settlement Name is required")
}

func TestCreateSettlement_ReportsCreationErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id").
		WithArgs(testUserID, "Fun time", 1, 0, 0, 1).
		WillReturnError(fmt.Errorf("insert error"))

	createRequest := settlement.CreateSettlementRequest{
		Name: "Fun time",
	}
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "return server error if creation fails")
}

func TestGetSettlement_ReturnsOneSettlement(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement where id = $1 AND owner = $2").
		WithArgs(1, testUserID).
		WillReturnRows(pgxmock.NewRows(settlementCols).
			AddRow(1, "owner", "Fun Forever", 1, 0, 0, 1))

	req := httptest.NewRequest("GET", "/settlements/1", nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "return OK on success")
	body, _ := io.ReadAll(resp.Body)
	dto := settlement.DTO{}
	_ = json.Unmarshal(body, &dto)
	assert.Equal(t, 1, dto.ID, "returned settlement should have supplied id")
}

func TestGetSettlement_ReportsScanErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement where id = $1 AND owner = $2").
		WithArgs(1, testUserID).
		WillReturnError(fmt.Errorf("scan error"))

	req := httptest.NewRequest("GET", "/settlements/1", nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "return server error on failure")
}

type WrongRequest struct {
	FancyName string `json:"soFancy"`
}
