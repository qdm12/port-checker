package server

import (
	"net"
	"net/http"
)

type Logger interface {
	Info(message string)
	Infof(format string, args ...any)
	Warn(message string)
	Errorf(format string, args ...any)
}

type RequestParser interface {
	ParseHTTPRequest(r *http.Request) net.IP
}
