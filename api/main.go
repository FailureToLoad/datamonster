package main

import (
	"context"
	"datamonster/settlement"
	"datamonster/user"
	userApi "datamonster/user/api"
	userRepo "datamonster/user/repo"
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
	appPool    *pgxpool.Pool
	piiPool    *pgxpool.Pool
	router     *chi.Mux
	appContext context.Context
)

func init() {
	appContext = context.Background()
	appPool = initAppPool(appContext)
	piiPool = initPrivatePool(appContext)
	router = web.NewRouter()
}

func main() {
	defer appPool.Close()

	settlementController := settlement.NewController(settlement.NewRepo(appPool))
	settlementController.RegisterRoutes(router)
	userController := userApi.NewController(user.NewService(userRepo.New(piiPool)))
	userController.RegisterRoutes(router)
	log.Default().Println("Starting server on port 8080")
	http.ListenAndServe(":8080", router)
}

func initAppPool(ctx context.Context) *pgxpool.Pool {
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

func initPrivatePool(ctx context.Context) *pgxpool.Pool {
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("PII_STRING"))
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
