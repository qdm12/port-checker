package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultListeningPort = "8000"

// Error codes:
// 1: Error in main server
// 2: Error in communicating with main server
// 3: Error in creating HTTP request
// 4: Error in parsing parameters

func main() {
	listeningPort := os.Getenv("PORT")
	if len(listeningPort) == 0 {
		listeningPort = defaultListeningPort
	} else {
		_, err := strconv.ParseInt(listeningPort, 10, 64)
		if err != nil {
			os.Exit(4)
		}
	}
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:"+listeningPort+"/healthcheck", nil)
	if err != nil {
		os.Exit(3)
	}
	client := &http.Client{Timeout: time.Duration(7000) * time.Millisecond}
	response, err := client.Do(request)
	if err != nil {
		os.Exit(2)
	}
	if response.StatusCode != 200 {
		os.Exit(1)
	}
}
