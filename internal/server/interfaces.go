package server

type Logger interface {
	Info(message string)
	Infof(format string, args ...any)
	Warn(message string)
	Errorf(format string, args ...any)
}
