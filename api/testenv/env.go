package testenv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/store/postgres/migrator"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBContainer struct {
	pgcontainer *postgres.PostgresContainer
	PGPool      *pgxpool.Pool
}

type Requester struct {
	authorizer *AuthorizerFake
	handler    http.Handler
}

func NewDBContainer(ctx context.Context) (*DBContainer, error) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %s", err)
	}

	dbsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %s", err)
	}

	pool, err := pgxpool.New(ctx, dbsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %s", err)
	}

	if err := migrator.Migrate(ctx, pool); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DBContainer{
		pgcontainer: postgresContainer,
		PGPool:      pool,
	}, nil
}

func (e *DBContainer) Cleanup() {
	e.PGPool.Close()
	if err := testcontainers.TerminateContainer(e.pgcontainer); err != nil {
		log.Fatalf("failed to terminate postgres container: %s", err)
	}
}

func NewRequester(controllers []server.Controller) (*Requester, error) {
	authorizer := &AuthorizerFake{
		authorized: true,
	}
	srv, err := server.New(AuthControllerFake{}, authorizer, []string{"localhost"}, controllers)
	if err != nil {
		return nil, err
	}

	return &Requester{
		handler:    srv.Handler,
		authorizer: authorizer,
	}, nil
}

func (r Requester) Unauthorized() func() {
	r.authorizer.SetAuthorized(false)

	return func() {
		r.authorizer.SetAuthorized(true)
	}
}

func (r Requester) DoRequest(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r Requester) CreateSettlement(userID string) (string, error) {
	body := `{"name":"Test Settlement"}`
	respBody, status := r.CreateSettlementWithBody(userID, body)
	if status != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d: %s", status, respBody.String())
	}

	var settlementID string
	err := json.NewDecoder(respBody).Decode(&settlementID)

	return settlementID, err
}

func (r Requester) CreateSettlementWithBody(userID, body string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.DoRequest(w, req)

	return w.Body, w.Code
}

func (r Requester) GetSettlements(userID string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()
	r.DoRequest(w, req)

	return w.Body, w.Code
}

func (r Requester) GetSettlement(userID string, settlementID string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID, nil)
	w := httptest.NewRecorder()
	r.DoRequest(w, req)

	return w.Body, w.Code
}

func (r Requester) CreateSurvivor(userID string, settlementID string, name string) (*bytes.Buffer, int) {
	body := fmt.Sprintf(`{"name":"%s","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`, name)
	return r.CreateSurvivorWithBody(userID, settlementID, body)
}

func (r Requester) CreateSurvivorWithBody(userID string, settlementID string, body string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.DoRequest(w, req)

	return w.Body, w.Code
}

func (r Requester) GetSurvivors(userID string, settlementID string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID+"/survivors", nil)
	w := httptest.NewRecorder()
	r.DoRequest(w, req)

	return w.Body, w.Code
}

func (r Requester) UpsertSurvivor(userID string, settlementID string, body string) (*bytes.Buffer, int) {
	r.authorizer.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodPut, "/api/settlements/"+settlementID+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.DoRequest(w, req)

	return w.Body, w.Code
}

func UUIDString() string {
	return UUID().String()
}

func UUID() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}
