package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

const defaultListeningPort = "80"

// HTMLData is used to contain data parsed for the Html template
type HTMLData struct {
	ClientIP string
	Browser  string
	Device   string
	OS       string
}

var indexPath = "index.html" // depends if in Docker with Scratch, we must hardcode it

func init() {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexPath = "/index.html"
	}
}

// Index is the main web UI response
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	browser, device, os := getUserAgentDetails(r.Header.Get("User-Agent"))
	clientIP, err := getClientIP(r)
	if err != nil {
		log.Println(err)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	data := HTMLData{
		clientIP,
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

func parseEnv() (listeningPort string, err error) {
	listeningPort = os.Getenv("PORT")
	if listeningPort == "" {
		listeningPort = defaultListeningPort
	} else {
		value, err := strconv.Atoi(listeningPort)
		if err != nil {
			return listeningPort, err
		}
		if value < 1 || value > 65535 {
			return listeningPort, errors.New("Environment variable PORT " + listeningPort + " is out of range (1 to 65535)")
		}
	}
	return listeningPort, nil
}

func main() {
	listeningPort, err := parseEnv()
	if err != nil {
		log.Fatalln(err)
	}
	router := httprouter.New()
	router.GET("/", Index)
	log.Println("Server started listening on port " + listeningPort)
	log.Fatal(http.ListenAndServe(":"+listeningPort, router))
}
