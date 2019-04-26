package server

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/julienschmidt/httprouter"
	"port-checker/pkg/logging"
	"port-checker/pkg/network"
)

type getParamsType struct {
	dir string
}

// CreateRouter returns a router with all the necessary routes configured
func CreateRouter(
	dir string,
) *httprouter.Router {
	getParams := getParamsType{
		dir: dir,
	}
	router := httprouter.New()
	router.POST("/", getParams.get)
	return router
}

func (params *getParamsType) get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	browser, device, os := network.GetUserAgentDetails(r.Header.Get("User-Agent"))
	ips, err := network.GetClientIP(r)
	if err != nil {
		logging.Warn("%s", err)
		return
	}
	logging.Success("received request from %s with device %s | %s | %s", ips.String(), device, os, browser)
	// HTMLData is used to contain data parsed for the Html template
	type HTMLData struct {
		ClientIP string
		Browser  string
		Device   string
		OS       string
	}
	data := HTMLData{
		ips.IP,
		browser,
		device,
		os,
	}
	t := template.Must(template.ParseFiles(params.dir + "/index.html"))
	err = t.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		logging.Warn("%s", err)
		fmt.Fprint(w, "Cannot create webpage: " + err.Error())
	}
}
