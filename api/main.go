package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/failuretoload/datamonster/auth"
	"github.com/failuretoload/datamonster/logger"
	"github.com/failuretoload/datamonster/server"
	"github.com/failuretoload/datamonster/settlement"
	settlementrepo "github.com/failuretoload/datamonster/settlement/repo"
	"github.com/failuretoload/datamonster/store/cache"
	"github.com/failuretoload/datamonster/store/postgres"
	"github.com/failuretoload/datamonster/store/postgres/migrator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger.Setup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sessions, err := cache.NewSessionStore(ctx)
	if err != nil {
		exit(fmt.Errorf("failed to initialize session store: %w", err))
	}
	authConfig := auth.Config{
		ClientID:      os.Getenv("CLIENT_ID"),
		ClientSecret:  os.Getenv("CLIENT_SECRET"),
		IssuerURL:     os.Getenv("ISSUER_URL"),
		RedirectURL:   os.Getenv("REDIRECT_URL"),
		IntrospectURL: os.Getenv("INTROSPECT_URL"),
		TokenURL:      os.Getenv("TOKEN_URL"),
		ClientURL:     os.Getenv("CLIENT_URL"),
		Sessions:      sessions,
	}

	authController, err := authConfig.Controller()
	if err != nil {
		exit(fmt.Errorf("failed to initialize auth controller: %w", err))
	}

	authorizer, err := authConfig.Authorizer()
	if err != nil {
		exit(fmt.Errorf("failed to initialize authorizer: %w", err))
	}

	dbsn := os.Getenv("DBSN")
	if dbsn == "" {
		exit(fmt.Errorf("dbsn is required"))
	}

	pool, err := postgres.NewConnectionPool(ctx, dbsn)
	if err != nil {
		exit(fmt.Errorf("failed to initialize connection pool: %w", err))
	}

	if err := migrator.Migrate(ctx, pool); err != nil {
		exit(fmt.Errorf("failed to run migrations: %w", err))
	}

	controllers, err := makeControllers(pool)
	if err != nil {
		exit(fmt.Errorf("failed to create controller: %w", err))
	}

	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		exit(errors.New("CLIENT_URL environment variable is required"))
	}

	srv, err := server.New(authController, authorizer, []string{clientURL}, controllers)
	if err != nil {
		exit(fmt.Errorf("failed to create server: %w", err))
	}

	go func() {
		slog.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		exit(fmt.Errorf("server forced to shutdown: %w", err))
	}

	slog.Info("server exited")
}

func exit(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

func makeControllers(pool *pgxpool.Pool) ([]server.Controller, error) {
	settlementRepo, err := settlementrepo.New(pool)
	if err != nil {
		return nil, err
	}
	settlementController, err := settlement.NewController(settlementRepo)
	if err != nil {
		return nil, err
	}

	return []server.Controller{
		settlementController,
	}, nil
}
