package main

import (
	"encoding/json"
	"errors"
	"regexp"
)

var regexIP = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`).FindString

func getSelfPublicIP() (pubIp string, err error) {
	content, err := getRequest("https://duckduckgo.com/?q=ip", httpGetTimeout)
	if err != nil {
		return pubIp, err
	}
	pubIp = regexIP(string(content))
	if pubIp == "" {
		return pubIp, errors.New("No public IP found on duckduck.go")
	}
	return pubIp, nil
}

type ipInfoType struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	GPS     string `json:"loc"`
	ISP     string `json:"org"`
}

func getLocationFromIP(ip string) (information *ipInfoType, err error) {
	content, err := getRequest("https://ipinfo.io/"+ip+"/json", httpGetTimeout)
	if err != nil {
		return nil, err
	}
	information = new(ipInfoType)
	err = json.Unmarshal(content, &information)
	if err != nil {
		return nil, err
	}
	return information, nil
}
