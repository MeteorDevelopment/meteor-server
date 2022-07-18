package api

import (
	"meteor-server/pkg/core"
	"net/http"
	"time"
)

var playing = make(map[string]time.Time)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	ip := core.IP(r)
	playing[ip] = time.Now()
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	ip := core.IP(r)

	delete(playing, ip)
}
