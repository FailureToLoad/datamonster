package settlement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"datamonster/settlement/mocks"
	"datamonster/settlement/repo"
	"datamonster/web"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

var (
	userId = 1
)

type SettlementApiTestSuite struct {
	suite.Suite
	target *Controller
	db     *mocks.MockConnection
	repo   *repo.PostgresRepo
	router *chi.Mux
}

func noopAuthHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), web.UserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (suite *SettlementApiTestSuite) SetupTest() {
	suite.db = &mocks.MockConnection{}
	suite.repo = repo.New(suite.db)
	suite.target = NewController(suite.repo)
	suite.router = chi.NewRouter()
	suite.target.RegisterRoutes(suite.router, noopAuthHandler)

}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ReturnsSettmentsList() {
	rows := mocks.MockRows{
		Rows: []pgx.Row{
			&mocks.SettlementRow{
				Id:                  1,
				Owner:               userId,
				Name:                "Fun Forever",
				SurvivalLimit:       1,
				DepartingSurvival:   0,
				CollectiveCognition: 0,
				CurrentYear:         1,
			},
			&mocks.SettlementRow{
				Id:                  2,
				Owner:               userId,
				Name:                "Wait, we get insanity for the croc?",
				SurvivalLimit:       1,
				DepartingSurvival:   0,
				CollectiveCognition: 0,
				CurrentYear:         1,
			},
		},
	}
	suite.db.SetRows(&rows)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	dto := SettlementsDTO{}
	json.Unmarshal(body, &dto)
	suite.Equal(2, dto.Count, "2 settlements should be returned")
}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ReportsScanErrors() {
	errorRows := mocks.MockRows{
		Rows: []pgx.Row{
			&mocks.SettlementRow{
				Id:                  1,
				Owner:               userId,
				Name:                "Fun Forever",
				SurvivalLimit:       1,
				DepartingSurvival:   0,
				CollectiveCognition: 0,
				CurrentYear:         1,
			},
			&mocks.ErrorRow{
				Error: fmt.Errorf("scan error"),
			},
		},
	}
	suite.db.SetRows(&errorRows)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "connection issues should result in server error")
}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ReportsConnectionErrors() {
	err := fmt.Errorf("query error")
	suite.db.SetError(err)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "deserialization issues should result in server error")
}

func (suite *SettlementApiTestSuite) Test_CreateSettlement_ReturnsASettlement() {
	insertRow := mocks.InsertRow{
		Id: 1,
	}
	suite.db.SetRow(&insertRow)
	settlementRequest := CreateSettlementRequest{
		Name: "Fun Forever",
	}
	reqBody, _ := json.Marshal(settlementRequest)
	req := httptest.NewRequest("POST", "/settlement", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(200, resp.StatusCode, "return 200 on success")
	respBody, _ := io.ReadAll(resp.Body)
	dto := SettlementDTO{}
	json.Unmarshal(respBody, &dto)
	suite.Equal(1, dto.Id, "created settlement should have an id")
	suite.Equal("Fun Forever", dto.Name, "created settlement should have supplied name")
}

type WrongRequest struct {
	FancyName string `json:"soFancy"`
}

func (suite *SettlementApiTestSuite) Test_CreateSettlement_EnforceRequestType() {
	wrongRequest := WrongRequest{
		FancyName: "Fun Forever",
	}
	reqBody, _ := json.Marshal(wrongRequest)
	req := httptest.NewRequest("POST", "/settlement", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(400, resp.StatusCode, "Request must be of type CreateSettlementRequest")
}

func (suite *SettlementApiTestSuite) Test_CreateSettlement_RequiresAName() {
	emptyRequest := CreateSettlementRequest{
		Name: "",
	}
	reqBody, _ := json.Marshal(emptyRequest)
	req := httptest.NewRequest("POST", "/settlement", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(400, resp.StatusCode, "Settlement Name is required")
}

func (suite *SettlementApiTestSuite) Test_CreateSettlement_ReportsCreationErrors() {
	insertRow := mocks.ErrorRow{
		Error: fmt.Errorf("insert error"),
	}
	suite.db.SetRow(&insertRow)
	createRequest := CreateSettlementRequest{
		Name: "Fun time",
	}
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/settlement", bytes.NewReader(reqBody))

	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "return server error if creation fails")
}

func (suite *SettlementApiTestSuite) Test_GetSettlement_ReturnsOneSettlement() {
	row := mocks.SettlementRow{
		Id:                  1,
		Owner:               userId,
		Name:                "Fun Forever",
		SurvivalLimit:       1,
		DepartingSurvival:   0,
		CollectiveCognition: 0,
		CurrentYear:         1,
	}
	suite.db.SetRow(&row)
	req := httptest.NewRequest("GET", "/settlement/1", nil)
	ctx := req.Context()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(200, resp.StatusCode, "return OK on success")
	body, _ := io.ReadAll(resp.Body)
	dto := SettlementDTO{}
	json.Unmarshal(body, &dto)
	suite.Equal(1, dto.Id, "returned settlement should have supplied id")
}

func (suite *SettlementApiTestSuite) Test_GetSettlement_ReportsScanErrors() {
	row := mocks.ErrorRow{
		Error: fmt.Errorf("scan error"),
	}
	suite.db.SetRow(&row)
	req := httptest.NewRequest("GET", "/settlement/1", nil)
	ctx := req.Context()
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, web.UserIdKey, userId)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "return server error on failure")
}
func TestSettlementApiTestSuite(t *testing.T) {
	suite.Run(t, new(SettlementApiTestSuite))
}
