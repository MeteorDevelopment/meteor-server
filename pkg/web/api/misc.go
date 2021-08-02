package api

import (
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
)

func SetDevBuildHandler(w http.ResponseWriter, r *http.Request) {
	devBuild := r.URL.Query().Get("devBuild")
	if devBuild != "" {
		db.SetDevBuild(devBuild)
	}

	core.Json(w, core.J{})
}
