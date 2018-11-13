package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/kyokomi/emoji"
)

const defaultListeningPort = "8000"

var fsLocation = ""

// HTMLData is used to contain data parsed for the Html template
type HTMLData struct {
	ClientIP string
	Browser  string
	Device   string
	OS       string
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	fsLocation = filepath.Dir(ex)
}

func parseEnv() (listeningPort string) {
	listeningPort = os.Getenv("PORT")
	if listeningPort == "" {
		listeningPort = defaultListeningPort
	} else {
		value, err := strconv.Atoi(listeningPort)
		if err != nil {
			log.Fatal(emoji.Sprint(":x:") + " PORT environment variable '" + listeningPort +
				"' is not a valid integer")
		}
		if value < 1024 {
			if os.Geteuid() == 0 {
				log.Println(emoji.Sprint(":warning:") + "PORT environment variable '" + listeningPort +
					"' allowed to be in the reserved system ports range as you are running as root.")
			} else if os.Geteuid() == -1 {
				log.Println(emoji.Sprint(":warning:") + "PORT environment variable '" + listeningPort +
					"' allowed to be in the reserved system ports range as you are running in Windows.")
			} else {
				log.Fatal(emoji.Sprint(":x:") + " PORT environment variable '" + listeningPort +
					"' can't be in the reserved system ports range (1 to 1023) when running without root.")
			}
		}
		if value > 65535 {
			log.Fatal(emoji.Sprint(":x:") + " PORT environment variable '" + listeningPort +
				"' can't be higher than 65535")
		}
		if value > 49151 {
			// dynamic and/or private ports.
			log.Println(emoji.Sprint(":warning:") + "PORT environment variable '" + listeningPort +
				"' is in the dynamic/private ports range (above 49151)")
		}
	}
	return listeningPort
}

func main() {
	fmt.Println("#################################")
	fmt.Println("######### Port Checker ##########")
	fmt.Println("######## by Quentin McGaw #######")
	fmt.Println("######## Give some " + emoji.Sprint(":heart:") + "at #########")
	fmt.Println("# github.com/qdm12/port-checker #")
	fmt.Print("#################################\n\n")
	listeningPort := parseEnv()
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/healthcheck", healthcheck)
	log.Println("Server listening on 0.0.0.0:" + listeningPort + emoji.Sprint(" :ear:"))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+listeningPort, router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	browser, device, os := getUserAgentDetails(r.Header.Get("User-Agent"))
	ips, err := getClientIP(r)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(emoji.Sprint(":heavy_check_mark:") + " received request from " + ips.String() +
		" with device " + device + "/" + os + "/" + browser)
	data := HTMLData{
		ips.ip,
		browser,
		device,
		os,
	}
	t := template.Must(template.ParseFiles(fsLocation + "/index.html"))
	err = t.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "An error occurred creating this webpage: "+err.Error())
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
