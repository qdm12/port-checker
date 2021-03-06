package server

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/avct/uasurfer"
	"github.com/qdm12/golibs/clientip"
	"github.com/qdm12/golibs/logging"
)

type handlers struct {
	// Config
	rootURL string
	// Objects
	logger        logging.Logger
	ipManager     clientip.Extractor
	indexTemplate *template.Template
	// Mockable functions
	timeNow func() time.Time
}

func newHandler(rootURL, templateStr string, logger logging.Logger,
	ipManager clientip.Extractor) (h http.Handler, err error) {
	indexTemplate, err := template.New("index.html").Parse(templateStr)
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
	ip := h.ipManager.HTTPRequest(r)

	h.logger.Info("received request from IP %s (device: %s | %s | %s)",
		ip, device, os, browser,
	)

	htmlData := struct {
		ClientIP string
		Browser  string
		Device   string
		OS       string
	}{
		ClientIP: ip.String(),
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
