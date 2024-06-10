package main

import (
	"context"
	"github.com/failuretoload/datamonster/settlement"
	postgres "github.com/failuretoload/datamonster/store/postgres"
	"github.com/failuretoload/datamonster/survivor"
	"github.com/failuretoload/datamonster/web"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	connPool   *pgxpool.Pool
	server     web.Server
	appContext context.Context
)

func init() {
	stInitErr := web.InitSuperTokens()
	if stInitErr != nil {
		log.Fatal(stInitErr)
	}
	appContext = context.Background()
	connPool = postgres.InitConnPool(appContext)
	server = web.NewServer()
}

func main() {
	defer connPool.Close()
	settlementController := settlement.NewController(connPool)
	survivorController := survivor.NewController(connPool)
	survivorController.RegisterRoutes(server.Mux)
	settlementController.RegisterRoutes(server.Mux)
	server.Start()
}
