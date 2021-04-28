package core

import (
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/wormhole/packets"
)

var clients map[ksuid.KSUID]*Client

func Broadcast(packet packets.PacketS2C) {
	for _, client := range clients {
		client.Send(packet)
	}
}

func AddClient(client *Client) {
	clients[client.ID] = client
}

func RemoveClient(client *Client) {
	delete(clients, client.ID)
}
