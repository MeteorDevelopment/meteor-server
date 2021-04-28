package wormhole

import (
	"github.com/gorilla/websocket"
	"meteor-server/pkg/wormhole/core"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := core.NewClient(conn)
	c.Start()

	c.Close()
}
