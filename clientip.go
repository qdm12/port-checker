package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var cidrs []*net.IPNet

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

func getClientIP(r *http.Request) (ip string, pubIP string, err error) {
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	var ipIsPrivate bool
	if xRealIP == "" && xForwardedFor == "" {
		ip = r.RemoteAddr
		if strings.ContainsRune(ip, ':') {
			ip, _, err = net.SplitHostPort(ip)
			if err != nil {
				return ip, pubIP, err
			}
		}
		ipIsPrivate, err = isPrivate(ip)
		if err != nil {
			return ip, pubIP, err
		}
		if ipIsPrivate {
			pubIP, err = getSelfPublicIP()
			if err != nil {
				return ip, pubIP, err
			}
		} else {
			pubIP = ip
		}
		return ip, pubIP, nil
	}
	for _, ip := range strings.Split(xForwardedFor, ",") {
		ip = strings.TrimSpace(ip)
		ipIsPrivate, err = isPrivate(ip)
		if err != nil {
			log.Println(err)
			continue
		}
		if !ipIsPrivate {
			return ip, pubIP, nil
		}
	}
	ip = xRealIP
	ipIsPrivate, err = isPrivate(ip)
	if err != nil {
		return ip, pubIP, err
	}
	if ipIsPrivate {
		pubIP, err = getSelfPublicIP()
		if err != nil {
			return ip, pubIP, err
		}
	} else {
		pubIP = ip
	}
	return ip, pubIP, nil
}

func getSelfPublicIP() (pubIp string, err error) {
	response, err := getRequest("https://ipinfo.io/ip", 1500)
	if err != nil {
		return pubIp, err
	}
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return pubIp, err
	}
	return string(content), nil
}
