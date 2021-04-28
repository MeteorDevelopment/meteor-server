package packets

import "encoding/json"

type PacketC2S struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PacketS2C struct {
	Type int         `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

const (
	Authenticate = iota
	JoinMessage  = iota
	LeaveMessage = iota
	Message      = iota
)
