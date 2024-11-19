package server

type Logger interface {
	Infof(format string, args ...any)
	Warn(message string)
	Errorf(format string, args ...any)
}
