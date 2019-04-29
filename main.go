package main

import (
	"fmt"
	"net/http"

	"github.com/kyokomi/emoji"

	"port-checker/pkg/healthcheck"
	"port-checker/pkg/logging"
	"port-checker/pkg/params"
	"port-checker/pkg/server"
)

func main() {
	if healthcheck.Mode() {
		healthcheck.Query()
	}
	fmt.Println("#################################")
	fmt.Println("######### Port Checker ##########")
	fmt.Println("######## by Quentin McGaw #######")
	fmt.Println("######## Give some " + emoji.Sprint(":heart:") + "at #########")
	fmt.Println("# github.com/qdm12/port-checker #")
	fmt.Print("#################################\n\n")
	logging.SetGlobalLoggerLevel(logging.InfoLevel)
	loggerMode := params.GetLoggerMode()
	logging.SetGlobalLoggerMode(loggerMode)
	nodeID := params.GetNodeID()
	logging.SetGlobalLoggerNodeID(nodeID)
	listeningPort := params.GetListeningPort()
	rootURL := params.GetRootURL()
	dir := params.GetDir()
	router := server.CreateRouter(rootURL, dir)
	logging.Info("Server listening on 0.0.0.0:%s%s", listeningPort, rootURL)
	err := http.ListenAndServe("0.0.0.0:"+listeningPort, router)
	if err != nil {
		logging.Fatal("%s", err)
	}
}
