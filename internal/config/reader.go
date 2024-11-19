package config

import (
	"github.com/qdm12/golibs/params"
)

type Reader struct {
	env params.Env
	os  params.OS
}

func NewReader() *Reader {
	return &Reader{
		env: params.NewEnv(),
		os:  params.NewOS(),
	}
}

func (r *Reader) ListeningPort() (port uint16, warning string, err error) {
	return r.env.ListeningPort("LISTENING_PORT", params.Default("8000"))
}

func (r *Reader) RootURL() (string, error) {
	return r.env.RootURL("ROOT_URL")
}
