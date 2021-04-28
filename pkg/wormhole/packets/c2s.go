package packets

type AuthenticateC2S struct {
	Token string `json:"token"`
}

type MessageC2S struct {
	Text string `json:"text"`
}
