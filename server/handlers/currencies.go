package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radekBednarik/cnb_forex/server/database"
)

func getCurrenciesSymbolsV1(c *gin.Context, dbs database.Database) {
	c.Header("Access-Control-Allow-Origin", "*")
	// query params
	now := time.Now()
	nowFormatted := now.Format("2006-01-02")
	weekTimeDelta := 7 * 24 * time.Hour
	dateWeekBeforeNow := now.Add(-weekTimeDelta)
	dwbnFormatted := dateWeekBeforeNow.Format("2006-01-02")
	dateFrom := c.DefaultQuery("dateFrom", dwbnFormatted)
	dateTo := c.DefaultQuery("dateTo", nowFormatted)

	data, err := dbs.SelectCurrenciesSymbolsV1(dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("DB error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetCurrenciesSymbolsV1(g *gin.Engine, dbs database.Database) {
	g.GET("/api/currencies/v1/symbols", func(c *gin.Context) {
		getCurrenciesSymbolsV1(c, dbs)
	})
}
