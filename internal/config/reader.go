package config

import (
	"github.com/qdm12/golibs/params"
)

type Reader interface {
	ListeningPort() (listeningPort, warning string, err error)
	RootURL() (string, error)
	ExeDir() (dir string, err error)
}

type reader struct {
	envParams params.EnvParams
}

func NewReader() Reader {
	return &reader{
		envParams: params.NewEnvParams(),
	}
}

func (r *reader) ListeningPort() (listeningPort, warning string, err error) {
	return r.envParams.GetListeningPort()
}

func (r *reader) RootURL() (string, error) {
	return r.envParams.GetRootURL()
}

func (r *reader) ExeDir() (dir string, err error) {
	return r.envParams.GetExeDir()
}
