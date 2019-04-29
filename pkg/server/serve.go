package server

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/julienschmidt/httprouter"
	"port-checker/pkg/logging"
	"port-checker/pkg/network"
)

type healthcheckParamsType struct {
}

type getParamsType struct {
	dir string
}

// CreateRouter returns a router with all the necessary routes configured
func CreateRouter(
	rootURL string,
	dir string,
) *httprouter.Router {
	healthcheckParams := healthcheckParamsType{}
	getParams := getParamsType{
		dir: dir,
	}
	router := httprouter.New()
	router.GET(rootURL+"/healthcheck", healthcheckParams.get)
	router.GET(rootURL+"/", getParams.get)
	return router
}

func (params *healthcheckParamsType) get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clientIP, err := network.GetClientIP(r)
	if err != nil {
		logging.Info("Cannot detect client IP: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if clientIP != "127.0.0.1" && clientIP != "::1" {
		logging.Info("IP address %s tried to perform the healthcheck", clientIP)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (params *getParamsType) get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	browser, device, os := network.GetUserAgentDetails(r.Header.Get("User-Agent"))
	headers := network.GetClientIPHeaders(r)
	ip, err := network.GetClientIP(r)
	if err != nil {
		logging.Warn("%s", err)
		return
	}
	logging.Success(
		"received request from IP %s (headers: %s) (device: %s | %s | %s)", 
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
	t := template.Must(template.ParseFiles(params.dir + "/index.html"))
	err = t.ExecuteTemplate(w, "index.html", htmlData)
	if err != nil {
		logging.Warn("%s", err)
		fmt.Fprint(w, "Cannot create webpage: " + err.Error())
	}
}
