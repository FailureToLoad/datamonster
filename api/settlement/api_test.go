package settlement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/failuretoload/datamonster/request"
	"github.com/pashagolub/pgxmock/v4"
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
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

func setupTest(t *testing.T) (*Controller, pgxmock.PgxPoolIface, *chi.Mux) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(sqlMatcher{}))
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	controller := NewController(mock)
	router := chi.NewRouter()
	controller.RegisterRoutes(router)

	return controller, mock, router
}

func TestGetSettlements_ReturnsSettlementsList(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	values := [][]any{
		{int32(1), testUserID, "Fun Forever", int32(1), int32(0), int32(0), int32(1)},
		{int32(2), testUserID, "Wait, we get insanity for the croc?", int32(1), int32(0), int32(0), int32(1)},
	}
	rows := pgxmock.NewRows(settlementCols).AddRows(values...)
	db.ExpectQuery("SELECT * FROM campaign.settlement WHERE owner = $1").
		WithArgs(testUserID).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/settlements", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req.WithContext(ctx))
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	var dto []DTO
	_ = json.Unmarshal(body, &dto)
	assert.Equal(t, 2, len(dto), "2 settlements should be returned")
}

func TestGetSettlements_ReportsScanErrors(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement WHERE owner=$1").
		WithArgs(testUserID).
		WillReturnError(fmt.Errorf("scan error"))

	req := httptest.NewRequest("GET", "/settlements", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "scan errors should result in server error")
}

func TestGetSettlements_ReportsConnectionErrors(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement WHERE owner=$1").
		WithArgs(testUserID).
		WillReturnError(fmt.Errorf("query error"))

	req := httptest.NewRequest("GET", "/settlements", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "connection issues should result in server error")
}

func TestCreateSettlement_ReturnsASettlement(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id").
		WithArgs(testUserID, "Fun Forever", int32(1), int32(0), int32(0), int32(1)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int32(1)))

	settlementRequest := CreateSettlementRequest{
		Name: "Fun Forever",
	}
	reqBody, _ := json.Marshal(settlementRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "return 200 on success")
	respBody, _ := io.ReadAll(resp.Body)
	dto := DTO{}
	_ = json.Unmarshal(respBody, &dto)
	assert.Equal(t, 1, dto.ID, "created settlement should have an id")
	assert.Equal(t, "Fun Forever", dto.Name, "created settlement should have supplied name")
}

func TestCreateSettlement_EnforceRequestType(t *testing.T) {
	_, _, router := setupTest(t)

	wrongRequest := WrongRequest{
		FancyName: "Fun Forever",
	}
	reqBody, _ := json.Marshal(wrongRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 400, resp.StatusCode, "Request must be of type CreateSettlementRequest")
}

func TestCreateSettlement_RequiresAName(t *testing.T) {
	_, _, router := setupTest(t)

	emptyRequest := CreateSettlementRequest{
		Name: "",
	}
	reqBody, _ := json.Marshal(emptyRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 400, resp.StatusCode, "Settlement Name is required")
}

func TestCreateSettlement_ReportsCreationErrors(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("INSERT INTO campaign.settlement (owner, name, survival_limit, departing_survival, collective_cognition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id").
		WithArgs(testUserID, "Fun time", int32(1), int32(0), int32(0), int32(1)).
		WillReturnError(fmt.Errorf("insert error"))

	createRequest := CreateSettlementRequest{
		Name: "Fun time",
	}
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/settlements", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "return server error if creation fails")
}

func TestGetSettlement_ReturnsOneSettlement(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement WHERE(id = $1 AND owner = $2)").
		WithArgs(int32(1), testUserID).
		WillReturnRows(pgxmock.NewRows(settlementCols).
			AddRow(int32(1), "owner", "Fun Forever", int32(1), int32(0), int32(0), int32(1)))

	req := httptest.NewRequest("GET", "/settlements/1", nil)
	ctx := req.Context()
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "return OK on success")
	body, _ := io.ReadAll(resp.Body)
	dto := DTO{}
	_ = json.Unmarshal(body, &dto)
	assert.Equal(t, 1, dto.ID, "returned settlement should have supplied id")
}

func TestGetSettlement_ReportsScanErrors(t *testing.T) {
	_, db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT * FROM campaign.settlement WHERE id = $1 AND owner = $2 LIMIT 1").
		WithArgs("1", testUserID).
		WillReturnError(fmt.Errorf("scan error"))

	req := httptest.NewRequest("GET", "/settlements/1", nil)
	ctx := req.Context()
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	ctx = context.WithValue(ctx, request.UserIDKey, testUserID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "return server error on failure")
}

type WrongRequest struct {
	FancyName string `json:"soFancy"`
}
