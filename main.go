package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qdm12/golibs/clientip"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosplash"
	"github.com/qdm12/log"
	"github.com/qdm12/port-checker/internal/config"
	"github.com/qdm12/port-checker/internal/server"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.New()

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, logger)
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
		if err == nil { // expected exit
			os.Exit(0)
		}
		logger.Warnf("Fatal error: %s", err)
		os.Exit(1)
	case signal := <-signalsCh:
		fmt.Println()
		logger.Warnf("Shutting down: signal %s", signal)
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
var templateStr string

var ErrPortOutOfRange = errors.New("port is out of range")

type Logger interface {
	Info(message string)
	Infof(format string, args ...any)
	Warn(message string)
	Errorf(format string, args ...any)
}

var (
	version string
	commit  string //nolint:gochecknoglobals
	created string //nolint:gochecknoglobals
)

func _main(ctx context.Context, logger Logger) (err error) {
	splashSettings := gosplash.Settings{
		User:       "qdm12",
		Repository: "port-checker",
		Emails:     []string{"quentin.mcgaw@gmail.com"},
		Version:    version,
		Commit:     commit,
		Created:    created,
		// Sponsor information
		PaypalUser:    "qmcgaw",
		GithubSponsor: "qdm12",
	}
	for _, line := range gosplash.MakeLines(splashSettings) {
		fmt.Println(line)
	}

	reader := reader.New(reader.Settings{})
	var settings config.Settings
	err = settings.Read(reader)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}
	settings.SetDefaults()
	err = settings.Validate()
	if err != nil {
		return fmt.Errorf("validating settings: %w", err)
	}
	fmt.Println(settings.String())

	ipManager := clientip.NewParser()

	server, err := server.New(*settings.ListeningAddress, *settings.RootURL, templateStr, logger, ipManager)
	if err != nil {
		return err
	}
	runError, err := server.Start(ctx)
	if err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}

	select {
	case err = <-runError:
		return err
	case <-ctx.Done():
		return server.Stop()
	}
}
