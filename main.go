package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type HtmlData struct {
	ClientIP    string
	ClientPubIP string
	City        string
	Region      string
	Country     string
	GPS         string
	ISP         string
	Browser     string
	Device      string
	OS          string
}

var indexPath = "index.html" // depends if in Docker with Scratch, we must hardcode it

func init() {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexPath = "/index.html"
	}
}

// TODO Google Maps with lat, long
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	browser, device, os := getUserAgentDetails(r.Header.Get("User-Agent"))
	clientIP, clientPubIP, err := getClientIP(r)
	if err != nil {
		log.Println(err)
		return
	}
	ipInfo, err := getLocationFromIP(clientPubIP)
	if err != nil {
		log.Println(err)
		return
	}
	data := HtmlData{
		clientIP,
		clientPubIP,
		ipInfo.City,
		ipInfo.Region,
		ipInfo.Country,
		ipInfo.GPS,
		ipInfo.ISP,
		browser,
		device,
		os,
	}
	t := template.Must(template.ParseFiles(indexPath))
	err = t.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	log.Fatal(http.ListenAndServe(":80", router))
}
