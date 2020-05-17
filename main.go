package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyokomi/emoji"

	"port-checker/pkg/handlers"
	"port-checker/pkg/params"

	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/server"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx))
}

func _main(ctx context.Context) int {
	if healthcheck.Mode(os.Args) {
		if err := healthcheck.Query(); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := createLogger(logging.InfoLevel)
	fmt.Println("#################################")
	fmt.Println("######### Port Checker ##########")
	fmt.Println("######## by Quentin McGaw #######")
	fmt.Println("######## Give some " + emoji.Sprint(":heart:") + "at #########")
	fmt.Println("# github.com/qdm12/port-checker #")
	fmt.Print("#################################\n\n")
	fatalOnError := func(err error) {
		if err != nil {
			logger.Error(err)
			cancel()
			time.Sleep(100 * time.Millisecond) // wait for operations to terminate
			os.Exit(1)
		}
	}
	paramsReader := params.NewReader()
	listeningPort, warning, err := paramsReader.GetListeningPort()
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	fatalOnError(err)
	rootURL, err := paramsReader.GetRootURL()
	fatalOnError(err)
	dir, err := paramsReader.GetDir()
	fatalOnError(err)
	ipManager := network.NewIPManager(logger)
	productionHandler := handlers.NewProductionHandler(rootURL, dir, ipManager, logger)
	healthcheckHandler := handlers.NewHealthcheckHandler()
	serverErrors := make(chan []error)
	go func() {
		logger.Info("Listening on port 0.0.0.0:%s at root URL: %s", listeningPort, rootURL)
		serverErrors <- server.RunServers(ctx,
			server.Settings{Name: "production", Addr: "0.0.0.0:" + listeningPort, Handler: productionHandler},
			server.Settings{Name: "healthcheck", Addr: "127.0.0.1:9999", Handler: healthcheckHandler},
		)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	select {
	case errors := <-serverErrors:
		for _, err := range errors {
			logger.Error(err)
		}
		return 1
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
		cancel()
		return 2
	case <-ctx.Done():
		logger.Warn("context canceled, shutting down")
		return 1
	}
}

func createLogger(level logging.Level) logging.Logger {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, level, -1)
	if err != nil {
		panic(err)
	}
	return logger
}
