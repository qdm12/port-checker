package server

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/avct/uasurfer"
	"github.com/qdm12/port-checker/internal/clientip"
)

type handlers struct {
	// Config
	listeningAddress string
	rootURL          string
	// Objects
	logger        Logger
	indexTemplate *template.Template
	// Mockable functions
	timeNow func() time.Time
}

func newHandler(listeningAddress, rootURL, templateStr string,
	logger Logger,
) (h http.Handler, err error) {
	indexTemplate, err := template.New("index.html").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	rootURL = strings.TrimRight(rootURL, "/")

	return &handlers{
		listeningAddress: listeningAddress,
		rootURL:          rootURL,
		logger:           logger,
		indexTemplate:    indexTemplate,
		timeNow:          time.Now,
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
	addrPort, err := clientip.ParseHTTPRequest(r)
	if err != nil {
		h.logger.Errorf("cannot parse ip address port: %s", err)
		http.Error(w, "cannot parse ip address port", http.StatusInternalServerError)
		return
	}

	var clientAddress string
	if addrPort.Port() == 0 {
		clientAddress = addrPort.Addr().String()
	} else {
		clientAddress = addrPort.String()
	}

	h.logger.Infof("received request from IP %s (device: %s | %s | %s)",
		clientAddress, device, os, browser,
	)

	htmlData := struct {
		ListeningAddress string
		ClientAddress    string
		Browser          string
		Device           string
		OS               string
	}{
		ListeningAddress: h.listeningAddress,
		ClientAddress:    clientAddress,
		Browser:          browser,
		Device:           device,
		OS:               os,
	}

	if err := h.indexTemplate.ExecuteTemplate(w, "index.html", htmlData); err != nil {
		h.logger.Errorf("executing template: %w", err)
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
