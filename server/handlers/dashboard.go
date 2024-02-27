package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radekBednarik/cnb_forex/server/database"
)

func getDashboardDataV1(c *gin.Context, dbs database.Database) {
	c.Header("Content-Type", "application/json")
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

	// stream data
	for name, dObject := range data.Data {
		toMarshall := make(map[string][]database.SingleCurrData)
		toMarshall[name] = dObject
		iJson, err := json.Marshal(toMarshall)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data returned from database."})
			return
		}
		c.Writer.Write(iJson)
		c.Writer.Flush()
	}

	c.Status(http.StatusOK)
}

func GetDashboardDataV1(s *gin.Engine, dbs database.Database) {
	s.GET("/api/dashboard/v1/data", func(c *gin.Context) {
		getDashboardDataV1(c, dbs)
	})
}
