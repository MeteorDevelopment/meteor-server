package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var playing = make(map[string]time.Time)
var uuids = make(map[string]string)

func PingHandler(c *gin.Context) {
	ip := c.ClientIP()
	playing[ip] = time.Now()

	id := c.Query("uuid")
	if id != "" {
		uuids[ip] = id
	}
}

func LeaveHandler(c *gin.Context) {
	ip := c.ClientIP()

	delete(playing, ip)
	delete(uuids, ip)
}

func UsingMeteorHandler(c *gin.Context) {
	var reqUuids []string
	err := c.BindJSON(&reqUuids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data."})
		return
	}

	resUuids := make(map[string]bool, len(reqUuids))
	for _, u := range uuids {
		for _, reqUuid := range reqUuids {
			if u == reqUuid {
				resUuids[u] = true
			}
		}
	}

	c.JSON(http.StatusOK, resUuids)
}

func ValidateOnlinePlayers() {
	now := time.Now()

	for ip, lastTimePlaying := range playing {
		if now.Sub(lastTimePlaying).Minutes() > 6 {
			delete(playing, ip)
			delete(uuids, ip)
		}
	}
}
