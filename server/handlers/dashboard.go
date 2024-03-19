package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radekBednarik/cnb_forex/server/database"
)

func getDashboardDataV1(c *gin.Context, dbs database.Database) {
	c.Header("Access-Control-Allow-Origin", "*")
	// query params
	now := time.Now()
	nowFormatted := now.Format("2006-01-02")
	weekTimeDelta := 7 * 24 * time.Hour
	dateWeekBeforeNow := now.Add(-weekTimeDelta)
	dwbnFormatted := dateWeekBeforeNow.Format("2006-01-02")
	dateFrom := c.DefaultQuery("dateFrom", dwbnFormatted)
	dateTo := c.DefaultQuery("dateTo", nowFormatted)

	// get data from database
	data, err := dbs.SelectDashboardDataV1(dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("DB error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, data)
}

func getDashboardDataV2(c *gin.Context, dbs database.Database) {
	c.Header("Access-Control-Allow-Origin", "*")
	// query params
	now := time.Now()
	nowFormatted := now.Format("2006-01-02")
	weekTimeDelta := 7 * 24 * time.Hour
	dateWeekBeforeNow := now.Add(-weekTimeDelta)
	dwbnFormatted := dateWeekBeforeNow.Format("2006-01-02")

	dateFrom := c.DefaultQuery("dateFrom", dwbnFormatted)
	dateTo := c.DefaultQuery("dateTo", nowFormatted)
	currency := c.DefaultQuery("currency", "USD")

	data, err := dbs.SelectDashboardDataV2(dateFrom, dateTo, currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("DB error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetDashboardDataV1(s *gin.Engine, dbs database.Database) {
	s.GET("/api/dashboard/v1/data", func(c *gin.Context) {
		getDashboardDataV1(c, dbs)
	})
}

func GetDashboardDataV2(s *gin.Engine, dbs database.Database) {
	s.GET("/api/dashboard/v2/data", func(c *gin.Context) {
		getDashboardDataV2(c, dbs)
	})
}
