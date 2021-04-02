package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"meteor-server/pkg/auth"
)

func LoginHandler(c *gin.Context) {
	token, err := auth.Login(c.Query("name"), c.Query("password"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong name or password."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LogoutHandler(c *gin.Context) {
	auth.Logout(c.GetString("id"))
}
