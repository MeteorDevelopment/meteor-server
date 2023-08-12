package core

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/db"
	"meteor-server/pkg/wormhole/packets"
)

type Client struct {
	conn  *websocket.Conn
	error bool

	Name string
	ID   ksuid.KSUID
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) IsAuthenticated() bool {
	return !c.ID.IsNil()
}

func (c *Client) Send(packet packets.PacketS2C) {
	err := c.conn.WriteJSON(packet)
	if err != nil {
		c.error = true
	}
}

func (c *Client) Start() {
	for {
		// Return from the loop if an error occurred
		if c.error {
			break
		}

		// Read packet
		var packet packets.PacketC2S
		err := c.conn.ReadJSON(&packet)
		if err != nil {
			break
		}

		if c.IsAuthenticated() {
			switch packet.Type {
			case packets.Message:
				var data packets.MessageC2S
				err = json.Unmarshal(packet.Data, &data)
				if err != nil {
					break
				}

				c.onMessage(data)
			}
		} else {
			// Authenticate
			if packet.Type == packets.Authenticate {
				var data packets.AuthenticateC2S
				err = json.Unmarshal(packet.Data, &data)
				if err != nil {
					break
				}

				c.onAuthenticate(data)
			}
		}
	}
}

func (c *Client) onAuthenticate(data packets.AuthenticateC2S) {
	// Validate token
	id, err := auth.IsTokenValid(data.Token)
	if err == nil {
		c.ID = id

		// Fetch name
		acc, err := db.GetAccountId(id)
		if err != nil {
			c.error = true
			return
		}

		c.Name = acc.Username

		// Send packet that confirms authentication
		c.Send(packets.NewAuthenticateS2C())

		// Send join message to current clients
		log.Info().Msgf("%s joined", c.Name)
		Broadcast(packets.NewJoinMessageS2C(c.Name))

		AddClient(c)
	}
}

func (c *Client) onMessage(data packets.MessageC2S) {
	Broadcast(packets.NewMessageS2C(c.Name, data.Text))
}

func (c *Client) Close() {
	_ = c.conn.Close()

	if c.IsAuthenticated() {
		RemoveClient(c)

		log.Info().Msgf("%s left", c.Name)
		Broadcast(packets.NewLeaveMessageS2C(c.Name))
	}
}
