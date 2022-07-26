package api

import (
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type apiAddonDeveloper struct {
	ID       ksuid.KSUID `json:"id"`
	Username string      `json:"username"`
}

type apiAddon struct {
	ID string `bson:"id" json:"id"`

	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Markdown    string `json:"markdown"`

	Developers []apiAddonDeveloper `json:"developers"`

	Version        string   `json:"version"`
	MeteorVersions []string `json:"meteor_versions"`
	Download       string   `json:"download"`

	DownloadCount int `json:"download_count"`

	Website string `json:"website"`
	Source  string `json:"source"`
	Support string `json:"support"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAddonById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		core.JsonError(w, "Invalid ID.")
		return
	}

	addon, err := db.GetAddon(id)
	if err != nil {
		core.JsonError(w, "No addon with this ID.")
		return
	}

	core.Json(w, getApiAddon(addon, true))
}

func SearchAddons(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	cursor, err := db.SearchAddons(r.URL.Query().Get("text"), page)
	//goland:noinspection GoUnhandledErrorResult
	defer cursor.Close(r.Context())

	if err != nil {
		core.JsonError(w, "Failed to search for addons.")
		return
	}

	addons := make([]apiAddon, 0)

	for {
		has := cursor.TryNext(r.Context())
		if !has {
			break
		}

		var addon db.Addon
		err := cursor.Decode(&addon)
		if err != nil {
			core.JsonError(w, "Failed to retrieve addons.")
			return
		}

		addons = append(addons, getApiAddon(addon, false))
	}

	core.Json(w, core.J{"addons": addons})
	_ = cursor.Close(nil)
}

func getApiAddon(addon db.Addon, includeMarkdown bool) apiAddon {
	markdown := ""
	if includeMarkdown {
		markdown = addon.Markdown
	}

	developers := make([]apiAddonDeveloper, 0, len(addon.Developers))

	for _, developer := range addon.Developers {
		account, err := db.GetAccountId(developer)

		if err == nil {
			developers = append(developers, apiAddonDeveloper{
				ID:       developer,
				Username: account.Username,
			})
		}
	}

	return apiAddon{
		ID:             addon.ID,
		Title:          addon.Title,
		Icon:           addon.Icon,
		Description:    addon.Description,
		Markdown:       markdown,
		Developers:     developers,
		Version:        addon.Version,
		MeteorVersions: addon.MeteorVersions,
		Download:       addon.Download,
		DownloadCount:  addon.DownloadCount,
		Website:        addon.Website,
		Source:         addon.Source,
		Support:        addon.Support,
		CreatedAt:      addon.CreatedAt,
		UpdatedAt:      addon.UpdatedAt,
	}
}
