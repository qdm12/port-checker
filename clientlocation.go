package main

import (
	"encoding/json"
	"io/ioutil"
)

type ipInfoType struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	GPS     string `json:"loc"`
	ISP     string `json:"org"`
}

func getLocationFromIP(ip string) (information *ipInfoType, err error) {
	response, err := getRequest("https://ipinfo.io/"+ip+"/json", 1500)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
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
