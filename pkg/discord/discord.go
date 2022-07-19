package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Guild       = "689197705683140636"
	MutedRole   = "741016178155192432"
	AccountRole = "777248653445300264"
	DonorRole   = "689205464574984353"
	DonorChat   = "713429344135020554"
	DevRole     = "689198253753106480"
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
	if !HasRole(user, role) {
		_ = send("PUT", "guilds/"+Guild+"/members/"+user+"/roles/"+role).Body.Close()
	}
}

func RemoveRole(user string, role string) {
	if HasRole(user, role) {
		_ = send("DELETE", "guilds/"+Guild+"/members/"+user+"/roles/"+role).Body.Close()
	}
}

func HasRole(user string, role string) bool {
	res := send("GET", "guilds/"+Guild+"/members/"+user)

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

func SendMessage(channel string, message string) {
	body := []byte(fmt.Sprintf(`{ "content": "%s" }`, message))
	req, _ := http.NewRequest("POST", "https://discord.com/api/channels/"+channel+"/messages", bytes.NewBuffer(body))
	req.Header.Set("User-Agent", "Meteor Server")
	req.Header.Set("Authorization", "Bot "+core.GetPrivateConfig().DiscordToken)
	req.Header.Set("Content-Type", "application/json")
	_, _ = client.Do(req)
}

func SendDonorMsg(user string) {
	SendMessage(DonorChat, fmt.Sprintf("<@%s> thanks for donating to Meteor and supporting the <@&%s> team.", user, DevRole))
}
