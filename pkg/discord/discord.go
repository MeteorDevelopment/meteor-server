package discord

import (
	"bytes"
	"encoding/json"
	"meteor-server/pkg/core"
	"net/http"
)

type User struct {
	Username      string
	Discriminator string
	Avatar        string
}

type member struct {
	Roles []string
}

const (
	guildId = "689197705683140636"

	MutedRole   = "741016178155192432"
	AccountRole = "777248653445300264"
	DonatorRole = "689205464574984353"
)

var client = http.Client{}

func send(method string, url string) *http.Response {
	req, _ := http.NewRequest(method, "https://discord.com/api/"+url, bytes.NewReader([]byte{}))
	req.Header.Set("User-Agent", "Meteor Server")
	req.Header.Set("Authorization", "Bot "+core.GetPrivateConfig().DiscordToken)

	res, _ := client.Do(req)
	return res
}

func GetUser(id string) User {
	res := send("GET", "users/"+id)

	var user User
	_ = json.NewDecoder(res.Body).Decode(&user)

	_ = res.Body.Close()
	return user
}

func AddRole(user string, role string) {
	_ = send("PUT", "guilds/"+guildId+"/members/"+user+"/roles/"+role).Body.Close()
}

func RemoveRole(user string, role string) {
	_ = send("DELETE", "guilds/"+guildId+"/members/"+user+"/roles/"+role).Body.Close()
}

func HasRole(user string, role string) bool {
	res := send("GET", "guilds/"+guildId+"/members/"+user)

	var member member
	_ = json.NewDecoder(res.Body).Decode(&member)

	_ = res.Body.Close()

	for _, r := range member.Roles {
		if r == role {
			return true
		}
	}

	return false
}
