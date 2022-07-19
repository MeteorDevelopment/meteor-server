package api

import (
	"meteor-server/pkg/db"
	"meteor-server/pkg/discord"
	"net/http"
)

func DiscordUserJoinedHandler(_ http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		account, err := db.GetAccountDiscordId(id)
		if err == nil {
			discord.AddRole(account.DiscordID, discord.AccountRole)

			if account.Donator {
				discord.AddRole(account.DiscordID, discord.DonorRole)
			}
		}
	}

	db.IncrementJoins()
}

func DiscordUserLeftHandler(_ http.ResponseWriter, _ *http.Request) {
	db.IncrementLeaves()
}
