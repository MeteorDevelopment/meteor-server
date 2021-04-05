package core

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GetDate() string {
	dt := time.Now()
	return fmt.Sprintf("%02d-%02d-%d", dt.Day(), dt.Month(), dt.Year())
}

func GetAccountID(c *gin.Context) ksuid.KSUID {
	id, _ := c.Get("id")
	return id.(ksuid.KSUID)
}

func IsEmailValid(email string) bool {
	length := len(email)
	if length < 3 && length > 254 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}
