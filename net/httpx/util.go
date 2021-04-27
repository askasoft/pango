package httpx

import (
	"net"
	"net/http"
	"strings"
)

// RemoteIPHeaders http remote ip headers
var RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-Ip"}

func validateIP(ip string) string {
	if ip == "" {
		return ""
	}

	i := strings.LastIndexByte(ip, ',')
	if i >= 0 {
		ip = ip[i+1:]
	}

	// // fix ipv4:port for IIS httpPlatformHandler
	i = strings.IndexByte(ip, ':')
	if i > 0 {
		ip = ip[0:i]
	}

	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ""
	}

	bs := net.ParseIP(ip)
	if bs == nil {
		return ""
	}

	return ip
}

// GetClientIP implements a best effort algorithm to return the real client IP.
// It will then try to parse the headers (defaulting to [X-Forwarded-For, X-Real-Ip]).
// If the headers are nots syntactically valid,
// the remote IP (coming form Request.RemoteAddr) is returned.
func GetClientIP(req *http.Request, headers ...string) string {
	if len(headers) == 0 {
		headers = RemoteIPHeaders
	}

	for _, name := range headers {
		ip := validateIP(req.Header.Get(name))
		if ip != "" {
			return ip
		}
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr))
	if err != nil {
		return ""
	}

	nip := net.ParseIP(ip)
	if nip == nil {
		return ""
	}
	return ip
}
