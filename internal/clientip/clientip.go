package clientip

import (
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
)

var ErrRequestIsNil = errors.New("request is nil")

func ParseHTTPRequest(request *http.Request) (ip netip.Addr, err error) {
	if request == nil {
		return ip, fmt.Errorf("%w", ErrRequestIsNil)
	}

	remoteAddress := removeSpaces(request.RemoteAddr)
	xRealIP := removeSpaces(request.Header.Get("X-Real-IP"))
	xForwardedFor := request.Header.Values("X-Forwarded-For")
	for i := range xForwardedFor {
		xForwardedFor[i] = removeSpaces(xForwardedFor[i])
	}

	// No header so it can only be remoteAddress
	if xRealIP == "" && len(xForwardedFor) == 0 {
		return getIPFromHostPort(remoteAddress)
	}

	// remoteAddress is the last proxy server forwarding the traffic
	// so we look into the HTTP headers to get the client IP
	xForwardedIPs := parseIPs(xForwardedFor)
	publicXForwardedIPs := extractPublicIPs(xForwardedIPs)
	if len(publicXForwardedIPs) > 0 {
		// first public XForwardedIP should be the client IP
		return publicXForwardedIPs[0], nil
	}

	// If all forwarded IP addresses are private we use the x-real-ip
	// address if it exists
	if xRealIP != "" {
		return getIPFromHostPort(xRealIP)
	}

	// Client IP is the first private IP address in the chain
	return xForwardedIPs[0], nil
}

func removeSpaces(header string) string {
	header = strings.ReplaceAll(header, " ", "")
	header = strings.ReplaceAll(header, "\t", "")
	return header
}
