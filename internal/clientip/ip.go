package clientip

import (
	"net/netip"
)

func extractPublicIPs(ips []netip.Addr) (publicIPs []netip.Addr) {
	publicIPs = make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		if ip.IsPrivate() {
			continue
		}
		publicIPs = append(publicIPs, ip)
	}
	return publicIPs
}

func parseIPs(stringIPs []string) (ips []netip.Addr) {
	for _, s := range stringIPs {
		ip, err := netip.ParseAddr(s)
		if err != nil {
			continue
		}
		ips = append(ips, ip)
	}
	return ips
}
