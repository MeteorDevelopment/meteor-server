package api

import (
	"meteor-server/pkg/db"
	"meteor-server/pkg/discord"
	"net/http"
)

func DiscordUserJoinedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		account, err := db.GetAccountDiscordId(id)
		if err == nil {
			discord.AddRole(account.DiscordID, discord.AccountRole)

			if account.Donator {
				discord.AddRole(account.DiscordID, discord.DonatorRole)
			}
		}
	}

	db.IncrementJoins()
}

func DiscordUserLeftHandler(w http.ResponseWriter, r *http.Request) {
	db.IncrementLeaves()
}

func GiveDonatorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		acc, err := db.GetAccountDiscordId(id)
		if err == nil {
			acc.GiveDonator()
		}
	}
}
