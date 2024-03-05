package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/radekBednarik/cnb_forex/server/database"
	"github.com/radekBednarik/cnb_forex/server/handlers"
)

func main() {
	g := gin.Default()

	// create db pool
	connString := fmt.Sprintf("user=%s password=%s host=localhost port=5432 dbname=cnb_forex sslmode=verify-ca pool_max_conns=16", os.Getenv("USER"), os.Getenv("PASSWORD"))
	dbs := database.Database{}
	dbs.New(connString)
	defer dbs.Pool.Close()

	// handle /api/dashboard/v1/data
	handlers.GetDashboardDataV1(g, dbs)
	// handle /api/currencies/v1/symbols
	handlers.GetCurrenciesSymbolsV1(g, dbs)

	g.Run()
}
