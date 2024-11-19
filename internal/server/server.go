// Package server runs the HTTP server for the program.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/golibs/clientip"
)

type Server struct {
	address string
	logger  Logger
	handler http.Handler
}

func New(address, rootURL, templateStr string,
	logger Logger, ipManager clientip.Extractor,
) (s *Server, err error) {
	handler, err := newHandler(rootURL, templateStr, logger, ipManager)
	if err != nil {
		return nil, err
	}
	return &Server{
		address: address,
		logger:  logger,
		handler: handler,
	}, nil
}

func (s *Server) Run(ctx context.Context, crashed chan<- error) {
	server := http.Server{Addr: s.address, Handler: s.handler, ReadHeaderTimeout: time.Second}
	go func() {
		<-ctx.Done()
		s.logger.Warn("shutting down (context canceled)")
		defer s.logger.Warn("shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil { //nolint:contextcheck
			s.logger.Errorf("failed shutting down: %s", err)
		}
	}()

	s.logger.Infof("listening on %s", s.address)
	crashed <- server.ListenAndServe()
}
