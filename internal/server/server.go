// Package server runs the HTTP server for the program.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/golibs/clientip"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, crashed chan<- error)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(address, rootURL, templateStr string,
	logger logging.Logger, ipManager clientip.Extractor,
) (s Server, err error) {
	handler, err := newHandler(rootURL, templateStr, logger, ipManager)
	if err != nil {
		return nil, err
	}
	return &server{
		address: address,
		logger:  logger,
		handler: handler,
	}, nil
}

func (s *server) Run(ctx context.Context, crashed chan<- error) {
	server := http.Server{Addr: s.address, Handler: s.handler, ReadHeaderTimeout: time.Second}
	go func() {
		<-ctx.Done()
		s.logger.Warn("shutting down (context canceled)")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()

	s.logger.Info("listening on %s", s.address)
	crashed <- server.ListenAndServe()
}
