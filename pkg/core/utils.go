package core

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

type J map[string]interface{}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func GetDate() string {
	dt := time.Now()
	return fmt.Sprintf("%02d-%02d-%d", dt.Day(), dt.Month(), dt.Year())
}

func GetAccountID(r *http.Request) ksuid.KSUID {
	id := r.Context().Value("id")
	return id.(ksuid.KSUID)
}

func IsEmailValid(email string) bool {
	length := len(email)
	if length < 3 && length > 254 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}

func IP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func Json(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

func JsonError(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(J{"error": message})
}
