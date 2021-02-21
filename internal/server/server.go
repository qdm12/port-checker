package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(ctx context.Context, address, rootURL, uiDir string,
	logger logging.Logger, ipManager network.IPManager) (s Server, err error) {
	handler, err := newHandler(rootURL, uiDir, logger, ipManager)
	if err != nil {
		return nil, err
	}
	return &server{
		address: address,
		logger:  logger,
		handler: handler,
	}, nil
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
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
	for ctx.Err() == nil {
		s.logger.Info("listening on %s", s.address)
		err := server.ListenAndServe()
		if err != nil && ctx.Err() == nil { // server crashed
			s.logger.Error(err)
			s.logger.Info("restarting")
		}
	}
}
