package config

import (
	"github.com/qdm12/golibs/params"
)

type Reader interface {
	ListeningPort() (port uint16, warning string, err error)
	RootURL() (string, error)
	ExeDir() (dir string, err error)
}

type reader struct {
	env params.Env
	os  params.OS
}

func NewReader() Reader {
	return &reader{
		env: params.NewEnv(),
		os:  params.NewOS(),
	}
}

func (r *reader) ListeningPort() (port uint16, warning string, err error) {
	return r.env.ListeningPort("LISTENING_PORT", params.Default("8000"))
}

func (r *reader) RootURL() (string, error) {
	return r.env.RootURL("ROOT_URL")
}

func (r *reader) ExeDir() (dir string, err error) {
	return r.os.ExeDir()
}
