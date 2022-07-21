package core

import (
	"encoding/json"
	"errors"
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

var cidrs []*net.IPNet

func Init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}
}

// IP https://github.com/tomasen/realip/blob/master/realip.go
func IP(r *http.Request) string {
	// Try Cf-Connecting-Ip
	cfConnectingIp := r.Header.Get("Cf-Connecting-Ip")

	if cfConnectingIp != "" {
		return cfConnectingIp
	}

	// Fetch header value
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	// If both empty, return IP from remote address
	if xRealIP == "" && xForwardedFor == "" {
		var remoteIP string

		// If there are colon in remote address, remove the port number
		// otherwise, return remote address as is
		if strings.ContainsRune(r.RemoteAddr, ':') {
			remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		} else {
			remoteIP = r.RemoteAddr
		}

		return remoteIP
	}

	// Check list of IP in X-Forwarded-For and return the first global address
	for _, address := range strings.Split(xForwardedFor, ",") {
		address = strings.TrimSpace(address)
		isPrivate, err := isPrivateAddress(address)
		if !isPrivate && err == nil {
			return address
		}
	}

	// If nothing succeed, return X-Real-IP
	return xRealIP
}

func isPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	for i := range cidrs {
		if cidrs[i].Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
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

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Failed to encode json response: ", err.Error())
		return
	}
}

func JsonError(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	err := json.NewEncoder(w).Encode(J{"error": message})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Failed to encode json error response: ", err.Error())
		return
	}
}

func Unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	err := json.NewEncoder(w).Encode(J{"error": "Unauthorized."})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Failed to encode json error response: ", err.Error())
		return
	}
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
