package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/avct/uasurfer"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

func NewHealthcheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewProductionHandler(rootURL string, exeDir string, ipManager network.IPManager, logger logging.Logger) http.HandlerFunc {
	logger = logger.WithPrefix("server: ")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || strings.TrimSuffix(r.URL.Path, "/") != strings.TrimSuffix(rootURL, "/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		browser, device, os := getUserAgentDetails(r.Header.Get("User-Agent"))
		headers := ipManager.GetClientIPHeaders(r)
		ip, err := ipManager.GetClientIP(r)
		if err != nil {
			logger.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Info("received request from IP %s (headers: %s) (device: %s | %s | %s)",
			ip, headers.String(), device, os, browser,
		)
		htmlData := struct {
			ClientIP string
			Browser  string
			Device   string
			OS       string
		}{
			ip,
			browser,
			device,
			os,
		}
		t := template.Must(template.ParseFiles(exeDir + "/index.html"))
		if err := t.ExecuteTemplate(w, "index.html", htmlData); err != nil {
			logger.Error(err)
			fmt.Fprint(w, "Cannot create webpage: "+err.Error())
		}
	}
}

func getUserAgentDetails(uaStr string) (browser, device, os string) {
	ua := uasurfer.Parse(uaStr)
	browser = fmt.Sprintf("%s %d", ua.Browser.Name.String()[7:], ua.Browser.Version.Major)
	device = ua.DeviceType.String()[6:]
	os = fmt.Sprintf("%s %d", ua.OS.Name.String()[2:], ua.OS.Version.Major)
	return browser, device, os
}
