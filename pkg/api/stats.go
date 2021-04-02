package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
)

type Stats struct {
	core.Config

	Date          string `json:"date"`
	DevBuild      string `json:"devBuild"`
	Downloads     int    `json:"downloads"`
	OnlinePlayers int    `json:"onlinePlayers"`
	OnlineUUIDs   int    `json:"onlineUUIDs"`
}

func StatsHandler(c *gin.Context) {
	date := c.Query("date")

	if date == "" {
		g := db.GetGlobal()

		c.JSON(http.StatusOK, Stats{Config: core.GetConfig(), Date: core.GetDate(), DevBuild: g.DevBuild, Downloads: g.Downloads, OnlinePlayers: len(playing), OnlineUUIDs: len(uuids)})
	} else {
		stats, err := db.GetJoinStats(date)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date."})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}
