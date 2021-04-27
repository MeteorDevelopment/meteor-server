package api

import (
	"net/http"

	"github.com/segmentio/ksuid"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	err := auth.Register(q.Get("username"), q.Get("email"), q.Get("password"))
	if err != nil {
		core.JsonError(w, err.Error())
		return
	}

	core.Json(w, core.J{})
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ksuid.Parse(r.URL.Query().Get("token"))
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	if !auth.ConfirmEmail(token) {
		core.JsonError(w, "Failed to confirm email address.")
		return
	}

	core.Json(w, core.J{})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	token, err := auth.Login(q.Get("name"), q.Get("password"))
	if err != nil {
		core.JsonError(w, "Wrong name or password.")
		return
	}

	core.Json(w, core.J{"token": token})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	auth.Logout(core.GetAccountID(r))
	core.Json(w, core.J{})
}

func AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	core.Json(w, account)
}

func McAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Get Minecraft UUID
	/*id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID."})
		return
	}

	// Get account
	account, err := db.GetAccount(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get account."})
		return
	}

	if c.Request.Method == "POST" { // Add Minecraft account
		err = account.AddMcAccount(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else { // Remove Minecraft account
		account.RemoveMcAccount(id)
	}

	c.JSON(http.StatusOK, gin.H{})*/
}
