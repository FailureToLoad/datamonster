package main

import (
	"context"
	"github.com/failuretoload/datamonster/settlement"
	settlementRepo "github.com/failuretoload/datamonster/settlement/repo"
	postgres "github.com/failuretoload/datamonster/store/postgres"
	"github.com/failuretoload/datamonster/survivor"
	"github.com/failuretoload/datamonster/web"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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
