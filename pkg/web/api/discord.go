package api

import (
	"meteor-server/pkg/db"
	"meteor-server/pkg/discord"
	"net/http"
)

func DiscordUserJoinedHandler(_ http.ResponseWriter, r *http.Request) {
	db.IncrementJoins()

	id := r.URL.Query().Get("id")
	if id == "" {
		return
	}

	account, err := db.GetAccountDiscordId(id)
	if err != nil {
		return
	}

	discord.AddRole(account.DiscordID, discord.AccountRole)

	if account.Donator {
		discord.AddRole(account.DiscordID, discord.DonorRole)
		discord.SendDonorMsg(account.DiscordID)
	}
}

func DiscordUserLeftHandler(_ http.ResponseWriter, _ *http.Request) {
	db.IncrementLeaves()
}
