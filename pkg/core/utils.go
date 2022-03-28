package core

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

type J map[string]interface{}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

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

func DownloadFile(formFile multipart.File, file *os.File, w http.ResponseWriter) bool {
	//goland:noinspection GoUnhandledErrorResult
	defer formFile.Close()
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	_, err := formFile.Seek(0, io.SeekStart)
	if err != nil {
		JsonError(w, "Server error. Failed to seek to the start of the file. Please contact developers.")
		return false
	}

	buf := make([]byte, 1024)
	for {
		n, err := formFile.Read(buf)
		if err != nil && err != io.EOF {
			JsonError(w, "Server error. Failed to read from sent cape file. Please contact developers.")
			return false
		}

		if n == 0 {
			break
		}

		_, err = file.Write(buf[:n])
		if err != nil {
			JsonError(w, "Server error. Failed to write to cape file. Please contact developers.")
			return false
		}
	}

	return true
}
