package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		id, err := IsTokenValid(token)

		if err == nil {
			c.Set("id", id)
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
	c.Abort()
}
