package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
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
			log.Println(err)
			os.Exit(1)
		}
		cidrs = append(cidrs, cidr)
	}
}

func isPrivate(address string) (bool, error) {
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

type ipHeaders struct {
	ip            string
	xRealIP       string
	xForwardedFor string
	remoteAddress string
}

func (ips *ipHeaders) String() string {
	return ips.ip + "(xRealIP=" + ips.xRealIP +
		", xForwardedFor=" + ips.xForwardedFor + ", remoteAddress=" +
		ips.remoteAddress + ")"
}

func getClientIP(r *http.Request) (ips ipHeaders, err error) {
	ips.xRealIP = r.Header.Get("X-Real-Ip")
	ips.xForwardedFor = r.Header.Get("X-Forwarded-For")
	ips.remoteAddress = r.RemoteAddr
	if ips.xRealIP == "" && ips.xForwardedFor == "" {
		ips.ip = ips.remoteAddress
		if strings.ContainsRune(ips.ip, ':') {
			ips.ip, _, err = net.SplitHostPort(ips.ip)
			if err != nil {
				return ips, err
			}
		}
		return ips, nil
	}
	for _, forwardedIP := range strings.Split(ips.xForwardedFor, ",") {
		ips.ip = strings.TrimSpace(forwardedIP)
		ipIsPrivate, err := isPrivate(ips.ip)
		if err != nil {
			log.Println(err)
			continue
		}
		if !ipIsPrivate {
			return ips, nil
		}
	}
	if ips.xRealIP == "" { // latest private xForwardedFor IP
		return ips, nil
	}
	return ips, nil
}
