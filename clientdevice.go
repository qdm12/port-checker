package main

import (
	"strconv"

	"github.com/avct/uasurfer"
)

func getUserAgentDetails(uaStr string) (browser, device, os string) {
	ua := uasurfer.Parse(uaStr)
	browser = ua.Browser.Name.String()[7:] + " " + strconv.FormatInt(int64(ua.Browser.Version.Major), 10)
	device = ua.DeviceType.String()[6:]
	os = ua.OS.Name.String()[2:] + " " + strconv.FormatInt(int64(ua.OS.Version.Major), 10)
	return browser, device, os
}
