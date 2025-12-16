package testenv

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/failuretoload/datamonster/server"
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

	schema := `
		CREATE SCHEMA campaign;

		CREATE TABLE campaign.settlement (
			id SERIAL PRIMARY KEY,
			external_id UUID NOT NULL DEFAULT gen_random_uuid(),
			owner VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			survival_limit INTEGER NOT NULL,
			departing_survival INTEGER NOT NULL,
			collective_cognition INTEGER NOT NULL,
			year INTEGER NOT NULL
		);

		CREATE INDEX idx_settlement_owner ON campaign.settlement(owner);
	`

	_, err = pool.Exec(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema: %w", err)
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
	r.authorizer.Authorized()
}

func (r Requester) Unauthorized() {
	r.authorizer.Unauthorized()
}

func (r Requester) ExpectUserID(userID string) {
	r.authorizer.ExpectUserID(userID)
}

func (r Requester) DoRequest(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
