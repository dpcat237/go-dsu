package httputil

import (
	"net"
	"time"
)

const connectionTimeout = 2 * time.Second

// IsConnection calls to Google's DNS server to check Internet connection
func IsConnection() bool {
	_, err := net.DialTimeout("tcp", "8.8.8.8:443", connectionTimeout)
	if err != nil {
		return false
	}
	return true
}
