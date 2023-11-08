package main

import (
	"context"
	"datamonster/settlement"
	"datamonster/web"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool       *pgxpool.Pool
	router     *chi.Mux
	appContext context.Context
)

func init() {
	appContext = context.Background()
	pool = initDbPool(appContext)
	router = web.NewRouter()
}

func main() {

	defer pool.Close()

	settlementController := settlement.NewController(settlement.NewRepo(pool))
	router.Route(settlement.BaseRoute, settlementController.RegisterRoutes)
	http.ListenAndServe(":8080", router)
}

func initDbPool(ctx context.Context) *pgxpool.Pool {
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to parse db config: %v\n", err)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, dbconfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return dbpool
}
