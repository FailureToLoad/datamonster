package main

import (
	"context"
	postgres "datamonster/connection/postgres"
	"datamonster/settlement"
	settlementRepo "datamonster/settlement/repo/postgres"
	"datamonster/user"
	userApi "datamonster/user/api"
	userRepo "datamonster/user/repo"
	"datamonster/web"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	appPool = postgres.InitAppPool(appContext)
	piiPool = postgres.InitPrivatePool(appContext)
	router = web.NewRouter()
}

func main() {
	defer appPool.Close()

	settlementController := settlement.NewController(settlementRepo.NewRepo(appPool))
	settlementController.RegisterRoutes(router)
	userController := userApi.NewController(user.NewService(userRepo.New(piiPool)))
	userController.RegisterRoutes(router)
	log.Default().Println("Starting server on port 8080")
	http.ListenAndServe(":8080", router)
}
