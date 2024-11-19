package clientip

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

func addrStringToAddrPort(address string) (addrPort netip.AddrPort, err error) {
	// address can be in the form ipv4:portStr, ipv6:portStr, ipv4 or ipv6
	host, portStr, err := splitHostPort(address)
	if err != nil {
		host = address
	}
	ip, err := netip.ParseAddr(host)
	if err != nil {
		return addrPort, fmt.Errorf("parsing IP address: %w", err)
	}
	var port uint16
	if portStr != "" {
		port, err = parsePort(portStr)
		if err != nil {
			return addrPort, fmt.Errorf("parsing port: %w", err)
		}
	}
	return netip.AddrPortFrom(ip, port), nil
}

func splitHostPort(address string) (ip, port string, err error) {
	if strings.ContainsRune(address, '[') && strings.ContainsRune(address, ']') {
		// should be an IPv6 address with brackets
		return net.SplitHostPort(address)
	}
	const ipv4MaxColons = 1
	if strings.Count(address, ":") > ipv4MaxColons {
		// could be an IPv6 without brackets
		i := strings.LastIndex(address, ":")
		port = address[i+1:]
		ip = address[0:i]
		if _, err := netip.ParseAddr(ip); err != nil {
			// invalid ip
			return net.SplitHostPort(address)
		}
		return ip, port, nil
	}
	// IPv4 address
	return net.SplitHostPort(address)
}

var ErrPortOutOfRange = errors.New("port is out of range")

func parsePort(s string) (port uint16, err error) {
	const base, bitSize = 10, 16
	portUint, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return 0, fmt.Errorf("parsing port: %w", err)
	} else if portUint > math.MaxUint16 {
		return 0, fmt.Errorf("%w: %d", ErrPortOutOfRange, portUint)
	}
	return uint16(portUint), nil
}
