package api

import (
	"encoding/json"
	"meteor-server/pkg/core"
	"net/http"
	"time"
)

var playing = make(map[string]time.Time)
var uuids = make(map[string]string)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	ip := core.IP(r)
	playing[ip] = time.Now()

	id := r.URL.Query().Get("uuid")
	if id != "" {
		uuids[ip] = id
	}
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	ip := core.IP(r)

	delete(playing, ip)
	delete(uuids, ip)
}

func UsingMeteorHandler(w http.ResponseWriter, r *http.Request) {
	var reqUuids []string
	err := json.NewDecoder(r.Body).Decode(&reqUuids)

	if err != nil {
		core.JsonError(w, "Invalid request data.")
		return
	}

	resUuids := make(map[string]bool, len(reqUuids))
	for _, u := range uuids {
		for _, reqUuid := range reqUuids {
			if u == reqUuid {
				resUuids[u] = true
			}
		}
	}

	core.Json(w, resUuids)
}

func ValidateOnlinePlayers() {
	now := time.Now()

	for ip, lastTimePlaying := range playing {
		if now.Sub(lastTimePlaying).Minutes() > 6 {
			delete(playing, ip)
			delete(uuids, ip)
		}
	}
}
