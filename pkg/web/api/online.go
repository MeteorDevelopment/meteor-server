package api

import (
	"meteor-server/pkg/core"
	"net/http"
	"sync"
	"time"
)

var mu = sync.RWMutex{}
var playing = make(map[string]time.Time)

func GetPlayingCount() int {
	mu.RLock()
	count := len(playing)
	mu.RUnlock()

	return count
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	ip := core.IP(r)
	playing[ip] = time.Now()

	mu.Unlock()
}

func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()

	ip := core.IP(r)
	delete(playing, ip)

	mu.Unlock()
}

func ValidateOnlinePlayers() {
	mu.Lock()
	now := time.Now()

	for ip, lastTimePlaying := range playing {
		if now.Sub(lastTimePlaying).Minutes() > 6 {
			delete(playing, ip)
		}
	}

	mu.Unlock()
}
