package clientip

import (
	"net"
	"net/netip"
	"strings"
)

func getIPFromHostPort(address string) (ip netip.Addr, err error) {
	// address can be in the form ipv4:port, ipv6:port, ipv4 or ipv6
	host, _, err := splitHostPort(address)
	if err != nil {
		host = address
	}
	return netip.ParseAddr(host)
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
