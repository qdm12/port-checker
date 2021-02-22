package server

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/avct/uasurfer"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type handlers struct {
	// Config
	rootURL string
	// Objects
	logger        logging.Logger
	ipManager     network.IPManager
	indexTemplate *template.Template
	// Mockable functions
	timeNow func() time.Time
}

func newHandler(rootURL, uiDir string, logger logging.Logger,
	ipManager network.IPManager) (h http.Handler, err error) {
	indexTemplate, err := parseIndexTemplate(uiDir)
	if err != nil {
		return nil, err
	}

	for strings.HasSuffix(rootURL, "/") {
		rootURL = strings.TrimSuffix(rootURL, "/")
	}

	return &handlers{
		rootURL:       rootURL,
		logger:        logger,
		ipManager:     ipManager,
		indexTemplate: indexTemplate,
		timeNow:       time.Now,
	}, nil
}

func (h *handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method != http.MethodGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case strings.TrimSuffix(r.URL.Path, "/") != h.rootURL:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	browser, device, os := getUserAgentDetails(r.Header.Get("User-Agent"))
	headers := h.ipManager.GetClientIPHeaders(r)
	ip, err := h.ipManager.GetClientIP(r)
	if err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("received request from IP %s (headers: %s) (device: %s | %s | %s)",
		ip, headers.String(), device, os, browser,
	)

	htmlData := struct {
		ClientIP string
		Browser  string
		Device   string
		OS       string
	}{
		ClientIP: ip,
		Browser:  browser,
		Device:   device,
		OS:       os,
	}

	if err := h.indexTemplate.ExecuteTemplate(w, "index.html", htmlData); err != nil {
		h.logger.Error(err)
		fmt.Fprint(w, "Cannot create webpage: "+err.Error())
	}
}

func getUserAgentDetails(uaStr string) (browser, device, os string) {
	ua := uasurfer.Parse(uaStr)
	browser = fmt.Sprintf("%s %d", ua.Browser.Name.String()[7:], ua.Browser.Version.Major)
	device = ua.DeviceType.String()[6:]
	os = fmt.Sprintf("%s %d", ua.OS.Name.String()[2:], ua.OS.Version.Major)
	return browser, device, os
}
