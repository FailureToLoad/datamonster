package settlement

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"datamonster/settlement/mocks"
	"datamonster/settlement/repo"
	"datamonster/web"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type SettlementApiTestSuite struct {
	suite.Suite
	target *Controller
	db     *mocks.MockConnection
	repo   *repo.PostgresRepo
}

func (suite *SettlementApiTestSuite) SetupTest() {
	suite.db = &mocks.MockConnection{}
	suite.repo = repo.New(suite.db)
	suite.target = NewController(suite.repo)
}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ShouldReturnSettmentsList() {
	settlementRows := []pgx.Row{
		&mocks.SettlementRow{
			Id:                  1,
			Owner:               1,
			Name:                "Fun Forever",
			SurvivalLimit:       1,
			DepartingSurvival:   0,
			CollectiveCognition: 0,
			CurrentYear:         1,
		},
		&mocks.SettlementRow{
			Id:                  2,
			Owner:               1,
			Name:                "Wait, we get insanity for the croc?",
			SurvivalLimit:       1,
			DepartingSurvival:   0,
			CollectiveCognition: 0,
			CurrentYear:         1,
		},
	}
	rows := mocks.MockRows{Rows: settlementRows}
	suite.db.SetRows(&rows)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, 1)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.target.getSettlements(w, req)
	resp := w.Result()

	suite.Equal(200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	dto := SettlementsDTO{}
	json.Unmarshal(body, &dto)
	suite.Equal(2, dto.Count, "2 settlements should be returned")
}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ShouldReportScanErrors() {
	errorRows := mocks.MockRows{
		Rows: []pgx.Row{
			&mocks.SettlementRow{
				Id:                  1,
				Owner:               1,
				Name:                "Fun Forever",
				SurvivalLimit:       1,
				DepartingSurvival:   0,
				CollectiveCognition: 0,
				CurrentYear:         1,
			},
			&mocks.ErrorRow{},
		},
	}
	suite.db.SetRows(&errorRows)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, 1)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.target.getSettlements(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "connection issues should result in server error")
}

func (suite *SettlementApiTestSuite) Test_GetSettlements_ShouldReportConnectionErrors() {
	err := fmt.Errorf("query error")
	suite.db.SetError(err)
	req := httptest.NewRequest("GET", "/settlement", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, web.UserIdKey, 1)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	suite.target.getSettlements(w, req)
	resp := w.Result()

	suite.Equal(500, resp.StatusCode, "connection issues should result in server error")
}

func TestSettlementApiTestSuite(t *testing.T) {
	suite.Run(t, new(SettlementApiTestSuite))
}
