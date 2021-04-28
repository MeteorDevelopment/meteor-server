package packets

// Authenticate

func NewAuthenticateS2C() PacketS2C {
	return PacketS2C{Type: Authenticate}
}

// Message

type MessageS2C struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func NewMessageS2C(user string, text string) PacketS2C {
	return PacketS2C{Type: Message, Data: MessageS2C{User: user, Text: text}}
}

// Join / Leave messages

type JoinLeaveS2C struct {
	User string `json:"user"`
}

func NewJoinMessageS2C(user string) PacketS2C {
	return PacketS2C{Type: JoinMessage, Data: JoinLeaveS2C{User: user}}
}

func NewLeaveMessageS2C(user string) PacketS2C {
	return PacketS2C{Type: LeaveMessage, Data: JoinLeaveS2C{User: user}}
}
