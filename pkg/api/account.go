package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/db"
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

func AccountInfoHandler(c *gin.Context) {
	account, err := db.GetAccount(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get account."})
		return
	}

	c.JSON(http.StatusOK, account)
}
