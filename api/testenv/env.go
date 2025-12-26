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
	"testing"

	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/store/postgres/migrator"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
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
	authorizer := &AuthorizerFake{}
	srv, err := server.New(AuthControllerFake{}, authorizer, []string{"localhost"}, controllers)
	if err != nil {
		return nil, err
	}

	return &Requester{
		handler:    srv.Handler,
		authorizer: authorizer,
	}, nil
}

func (r Requester) Authorized() {
	r.authorizer.SetAuthorized(true)
}

func (r Requester) Unauthorized() {
	r.authorizer.SetAuthorized(false)
}

func (r Requester) ExpectUserID(userID string) {
	r.authorizer.ExpectUserID(userID)
}

func (r Requester) DoRequest(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r Requester) CreateSettlement(t *testing.T, userID string) uuid.UUID {
	return r.CreateSettlementWithName(t, userID, "Test Settlement")
}

func (r Requester) CreateSettlementWithName(t *testing.T, userID, name string) uuid.UUID {
	r.Authorized()
	r.ExpectUserID(userID)

	body := fmt.Sprintf(`{"name":"%s"}`, name)
	req := httptest.NewRequest(http.MethodPost, "/api/settlements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.DoRequest(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var settlementID uuid.UUID
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settlementID))
	return settlementID
}

func (r Requester) GetSettlements(t *testing.T, userID string) *bytes.Buffer {
	r.Authorized()
	r.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements", nil)
	w := httptest.NewRecorder()

	r.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	return w.Body
}

func (r Requester) GetSettlement(t *testing.T, userID string, settlementID uuid.UUID) *bytes.Buffer {
	require.NotEqual(t, uuid.Nil, settlementID)

	r.Authorized()
	r.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID.String(), nil)
	w := httptest.NewRecorder()

	r.DoRequest(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	return w.Body
}

func (r Requester) CreateSurvivor(t *testing.T, userID string, settlementID uuid.UUID, name string) *bytes.Buffer {
	r.Authorized()
	r.ExpectUserID(userID)

	body := fmt.Sprintf(`{"name":"%s","birth":1,"gender":"M","huntxp":0,"survival":1,"movement":5,"accuracy":0,"strength":0,"evasion":0,"luck":0,"speed":0,"insanity":0,"systemicPressure":0,"torment":0,"lumi":0,"courage":0,"understanding":0}`, name)
	req := httptest.NewRequest(http.MethodPost, "/api/settlements/"+settlementID.String()+"/survivors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.DoRequest(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	return w.Body
}

func (r Requester) GetSurvivors(t *testing.T, userID string, settlementID uuid.UUID) *bytes.Buffer {
	r.Authorized()
	r.ExpectUserID(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/settlements/"+settlementID.String()+"/survivors", nil)
	w := httptest.NewRecorder()

	r.DoRequest(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	return w.Body
}
