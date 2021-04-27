package api

import (
	"net/http"

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

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")

	if date == "" {
		g := db.GetGlobal()

		core.Json(w, Stats{Config: core.GetConfig(), Date: core.GetDate(), DevBuild: g.DevBuild, Downloads: g.Downloads, OnlinePlayers: len(playing), OnlineUUIDs: len(uuids)})
	} else {
		stats, err := db.GetJoinStats(date)

		if err != nil {
			core.JsonError(w, "Invalid date.")
			return
		}

		core.Json(w, stats)
	}
}
