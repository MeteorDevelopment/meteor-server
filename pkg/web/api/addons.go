package api

import (
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
)

func GetAddonById(w http.ResponseWriter, r *http.Request) {
	id, err := ksuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		core.JsonError(w, "Invalid ID.")
		return
	}

	addon, err := db.GetAddon(id)
	if err != nil {
		core.JsonError(w, "No addon with this ID.")
		return
	}

	core.Json(w, addon)
}

func SearchAddons(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.SearchAddons(r.URL.Query().Get("text"))
	if err != nil {
		core.JsonError(w, "Failed to search for addons.")
		return
	}

	addons := make([]db.Addon, 0)
	_ = cursor.All(nil, &addons)

	core.Json(w, core.J{"addons": addons})
	_ = cursor.Close(nil)
}
