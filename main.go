package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"port-checker/internal/config"
	"port-checker/internal/health"
	"port-checker/internal/server"
	"strconv"
	"syscall"
	"time"

	"github.com/qdm12/golibs/clientip"
	"github.com/qdm12/golibs/logging"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := logging.New(logging.StdLog)

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, os.Args, logger)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	select {
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Warn("Fatal error:", err)
		os.Exit(1)
	case signal := <-signalsCh:
		fmt.Println()
		logger.Warn("Shutting down: signal", signal)
	}

	cancel()

	const shutdownGracePeriod = time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case <-errorCh:
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
		logger.Warn("Shutdown timed out")
	}

	os.Exit(1)
}

//go:embed index.html
var templateStr string //nolint:gochecknoglobals

func _main(ctx context.Context, args []string, logger logging.Logger) error {
	if health.IsClientMode(args) {
		client := health.NewClient()
		if err := client.Query(ctx); err != nil {
			return err
		}
		return nil
	}

	fmt.Println("#################################")
	fmt.Println("######### Port Checker ##########")
	fmt.Println("######## by Quentin McGaw #######")
	fmt.Println("######## Give some ❤️ at #########")
	fmt.Println("# github.com/qdm12/port-checker #")
	fmt.Print("#################################\n\n")
	paramsReader := config.NewReader()
	listeningPort, warning, err := paramsReader.ListeningPort()
	if len(warning) > 0 {
		logger.Warn(warning)
	}
	if err != nil {
		return err
	}

	rootURL, err := paramsReader.RootURL()
	if err != nil {
		return err
	}

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	healthServer := flagSet.Bool("healthserver", false, "Enable the health HTTP server")
	listeningPortPtr := flagSet.String("port", "", "Listening port for the HTTP server")
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	if portStr := *listeningPortPtr; len(portStr) > 0 {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		listeningPort = uint16(port)
	}

	ipManager := clientip.NewExtractor()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	N := 0
	crashed := make(chan error)

	if *healthServer {
		healthcheckServer := health.NewServer("127.0.0.1:9999",
			logger.NewChild(logging.SetPrefix("healthcheck: ")), health.MakeIsHealthy())
		N++
		go healthcheckServer.Run(ctx, crashed)
	}

	address := "0.0.0.0:" + strconv.FormatInt(int64(listeningPort), 10)
	server, err := server.New(address, rootURL, templateStr, logger, ipManager)
	if err != nil {
		return err
	}
	N++
	go server.Run(ctx, crashed)

	err = <-crashed
	cancel()

	for i := 1; i < N; i++ {
		<-crashed
	}

	return err
}
