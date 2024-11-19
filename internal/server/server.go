// Package server runs the HTTP server for the program.
package server

import (
	"github.com/qdm12/goservices/httpserver"
)

func New(address, rootURL, templateStr string,
	logger Logger, ipManager RequestParser,
) (s *httpserver.Server, err error) {
	handler, err := newHandler(rootURL, templateStr, logger, ipManager)
	if err != nil {
		return nil, err
	}
	return httpserver.New(httpserver.Settings{
		Handler: handler,
		Address: &address,
		Logger:  logger,
	})
}
