package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"port-checker/internal/config"
	"port-checker/internal/health"
	"port-checker/internal/server"
	"sync"
	"syscall"

	"github.com/kyokomi/emoji"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx))
}

func _main(ctx context.Context) int {
	if health.IsClientMode(os.Args) {
		client := health.NewClient()
		if err := client.Query(ctx); err != nil {
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
	paramsReader := config.NewReader()
	listeningPort, warning, err := paramsReader.ListeningPort()
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	if err != nil {
		logger.Error(err)
		return 1
	}

	rootURL, err := paramsReader.RootURL()
	if err != nil {
		logger.Error(err)
		return 1
	}

	dir, err := paramsReader.ExeDir()
	if err != nil {
		logger.Error(err)
		return 1
	}

	ipManager := network.NewIPManager(logger)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	healthcheckServer := health.NewServer("127.0.0.1:9999",
		logger.WithPrefix("healthcheck: "), health.MakeIsHealthy())
	wg.Add(1)
	go healthcheckServer.Run(ctx, wg)

	server, err := server.New(ctx, "0.0.0.0:"+listeningPort,
		rootURL, dir, logger, ipManager)
	if err != nil {
		logger.Error(err)
		return 1
	}
	wg.Add(1)
	go server.Run(ctx, wg)

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
		cancel()
	case <-ctx.Done():
		logger.Warn("context canceled, shutting down")
	}
	return 1
}

func createLogger(level logging.Level) logging.Logger {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, level, -1)
	if err != nil {
		panic(err)
	}
	return logger
}
