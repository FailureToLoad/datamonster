package main

import (
	"context"
	"datamonster/settlement"
	settlementRepo "datamonster/settlement/repo"
	postgres "datamonster/store/postgres"
	"datamonster/survivor"
	"datamonster/web"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var (
	appPool    *pgxpool.Pool
	server     web.Server
	appContext context.Context
)

func init() {
	stInitErr := web.InitSuperTokens()
	if stInitErr != nil {
		log.Fatal(stInitErr)
	}
	appContext = context.Background()
	appPool = postgres.InitAppPool(appContext)
	server = web.NewServer()
}

func main() {
	defer appPool.Close()
	settlementController := settlement.NewController(settlementRepo.New(appPool))
	survivorController := survivor.NewController(appPool)
	survivorController.RegisterRoutes(server.Mux)
	settlementController.RegisterRoutes(server.Mux)
	server.Start()
}
