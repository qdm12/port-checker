package network

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"port-checker/pkg/logging"
)

var cidrs []*net.IPNet
var regexIP = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`).FindString

func init() {
	maxCidrBlocks := [8]string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for _, maxCidrBlock := range maxCidrBlocks {
		_, cidr, err := net.ParseCIDR(maxCidrBlock)
		if err != nil {
			logging.Fatal("%s", err)
		}
		cidrs = append(cidrs, cidr)
	}
}

func ipIsPrivate(ip string) (bool, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false, fmt.Errorf("address %s is not valid", ip)
	}
	for i := range cidrs {
		if cidrs[i].Contains(netIP) {
			return true, nil
		}
	}
	return false, nil
}

// IPHeaders contains all the IP related headers of an HTTP request
type IPHeaders struct {
	IP            string
	XRealIP       string
	XForwardedFor string
	RemoteAddress string
}

func (ips *IPHeaders) String() string {
	return fmt.Sprintf("%s (XRealIP=%s | XForwardedFor=%s | RemoteAddress=%s)", ips.IP, ips.XRealIP, ips.XForwardedFor, ips.RemoteAddress)
}

// GetClientIP returns all the IP addresses found in a HTTP request
func GetClientIP(r *http.Request) (ips IPHeaders, err error) {
	ips.XRealIP = r.Header.Get("X-Real-Ip")
	ips.XForwardedFor = r.Header.Get("X-Forwarded-For")
	ips.RemoteAddress = r.RemoteAddr
	if ips.XRealIP == "" && ips.XForwardedFor == "" {
		ips.IP = ips.RemoteAddress
		if strings.ContainsRune(ips.IP, ':') {
			ips.IP, _, err = net.SplitHostPort(ips.IP)
			if err != nil {
				return ips, err
			}
		}
		return ips, nil
	}
	for _, forwardedIP := range strings.Split(ips.XForwardedFor, ",") {
		ips.IP = strings.TrimSpace(forwardedIP)
		private, err := ipIsPrivate(ips.IP)
		if err != nil {
			logging.Warn("%s", err)
			continue
		}
		if !private {
			return ips, nil
		}
	}
	if ips.XRealIP == "" { // latest private XForwardedFor IP
		return ips, nil
	}
	return ips, nil
}
