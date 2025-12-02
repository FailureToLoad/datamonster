package survivor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/failuretoload/datamonster/request"
	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testSettlementID = 1
	testUserID       = "userId"
	survivorCols     = []string{
		"id", "settlement", "name", "birth", "gender",
		"status", "huntxp", "survival", "movement",
		"accuracy", "strength", "evasion", "luck",
		"speed", "insanity", "systemic_pressure",
		"torment", "lumi", "courage", "understanding",
	}
)

func setupTest(t *testing.T) (pgxmock.PgxPoolIface, *chi.Mux) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	controller := NewController(mock)
	router := chi.NewRouter()
	router.Use(request.AutheliaMiddleware)
	controller.RegisterRoutes(router)

	return mock, router
}

func TestGetSurvivors_ReturnsSurvivorList(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	values := [][]any{
		{1, testSettlementID, "Survivor One", 1, "F", nil, 0, 1, 5, 5, 5, 5, 5, 5, 0, 0, 0, 0, 0, 0},
		{2, testSettlementID, "Survivor Two", 1, "M", nil, 0, 1, 5, 5, 5, 5, 5, 5, 0, 0, 0, 0, 0, 0},
	}
	rows := pgxmock.NewRows(survivorCols).AddRows(values...)
	db.ExpectQuery("SELECT .* FROM campaign.survivor WHERE").
		WithArgs(testSettlementID).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	var dto []DTO
	_ = json.Unmarshal(body, &dto)
	assert.Equal(t, 2, len(dto), "2 survivors should be returned")
}

func TestGetSurvivors_ReportsScanErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT .* FROM campaign.survivor WHERE").
		WithArgs(testSettlementID).
		WillReturnError(fmt.Errorf("scan error"))

	req := httptest.NewRequest("GET", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "scan errors should result in server error")
}

func TestGetSurvivors_ReportsConnectionErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectQuery("SELECT .* FROM campaign.survivor WHERE").
		WithArgs(testSettlementID).
		WillReturnError(fmt.Errorf("connection error"))

	req := httptest.NewRequest("GET", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), nil)
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "connection issues should result in server error")
}

func TestCreateSurvivor_Success(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectExec("INSERT INTO campaign.survivor").
		WithArgs(
			testSettlementID, "New Survivor", 1, "F",
			0, 1, 5, 5, 5,
			5, 5, 5, 0, 0,
			0, 0, 0, 0,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	survivorRequest := DTO{
		Settlement: testSettlementID,
		Name:       "New Survivor",
		Birth:      1,
		Gender:     "F",
		HuntXp:     0,
		Survival:   1,
		Movement:   5,
		Accuracy:   5,
		Strength:   5,
		Evasion:    5,
		Luck:       5,
		Speed:      5,
	}
	reqBody, _ := json.Marshal(survivorRequest)
	req := httptest.NewRequest("POST", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 204, resp.StatusCode, "successful creation should return 204")
}

func TestCreateSurvivor_RequiresValidSettlementID(t *testing.T) {
	_, router := setupTest(t)

	survivorRequest := DTO{
		Name:   "New Survivor",
		Gender: "F",
	}
	reqBody, _ := json.Marshal(survivorRequest)
	req := httptest.NewRequest("POST", "/settlements/invalid/survivors", bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "invalid settlement ID should return 500")
}

func TestCreateSurvivor_RequiresName(t *testing.T) {
	_, router := setupTest(t)

	survivorRequest := DTO{
		Settlement: testSettlementID,
		Gender:     "F",
	}
	reqBody, _ := json.Marshal(survivorRequest)
	req := httptest.NewRequest("POST", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "missing name should return 500")
}

func TestCreateSurvivor_RequiresUniqueName(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectExec("INSERT INTO campaign.survivor").
		WithArgs(
			testSettlementID, "Duplicate Name", 1, "F",
			0, 1, 5, 5, 5,
			5, 5, 5, 0, 0,
			0, 0, 0, 0,
		).
		WillReturnError(fmt.Errorf("UNIQUE constraint failed"))

	survivorRequest := DTO{
		Settlement: testSettlementID,
		Name:       "Duplicate Name",
		Birth:      1,
		Gender:     "F",
		HuntXp:     0,
		Survival:   1,
		Movement:   5,
		Accuracy:   5,
		Strength:   5,
		Evasion:    5,
		Luck:       5,
		Speed:      5,
	}
	reqBody, _ := json.Marshal(survivorRequest)
	req := httptest.NewRequest("POST", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 400, resp.StatusCode, "duplicate name should return 400")
}

func TestCreateSurvivor_ReportsCreationErrors(t *testing.T) {
	db, router := setupTest(t)
	defer db.Close()

	db.ExpectExec("INSERT INTO campaign.survivor").
		WithArgs(
			testSettlementID, "Error Survivor", 1, "F",
			0, 1, 5, 5, 5,
			5, 5, 5, 0, 0,
			0, 0, 0, 0,
		).
		WillReturnError(fmt.Errorf("database error"))

	survivorRequest := DTO{
		Settlement: testSettlementID,
		Name:       "Error Survivor",
		Birth:      1,
		Gender:     "F",
		HuntXp:     0,
		Survival:   1,
		Movement:   5,
		Accuracy:   5,
		Strength:   5,
		Evasion:    5,
		Luck:       5,
		Speed:      5,
	}
	reqBody, _ := json.Marshal(survivorRequest)
	req := httptest.NewRequest("POST", fmt.Sprintf("/settlements/%d/survivors", testSettlementID), bytes.NewReader(reqBody))
	req.Header.Set("Remote-User", testUserID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, 500, resp.StatusCode, "database errors should return 500")
}
