package core

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

func GetDate() string {
	dt := time.Now()
	return fmt.Sprintf("%02d-%02d-%d", dt.Day(), dt.Month(), dt.Year())
}

func GetAccountID(c *gin.Context) ksuid.KSUID {
	id, _ := c.Get("id")
	return id.(ksuid.KSUID)
}
