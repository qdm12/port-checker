package params

import (
	libparams "github.com/qdm12/golibs/params"
)

type Reader interface {
	GetListeningPort() (listeningPort, warning string, err error)
	GetRootURL() (string, error)
	GetDir() (dir string, err error)
}

type reader struct {
	envParams libparams.EnvParams
}

func NewReader() Reader {
	return &reader{
		envParams: libparams.NewEnvParams(),
	}
}

// GetListeningPort obtains and checks the listening port from Viper (env variable or config file, etc.)
func (r *reader) GetListeningPort() (listeningPort, warning string, err error) {
	return r.envParams.GetListeningPort()
}

// GetRootURL obtains and checks the root URL from Viper (env variable or config file, etc.)
func (r *reader) GetRootURL() (string, error) {
	return r.envParams.GetRootURL()
}

// GetDir obtains the executable directory
func (r *reader) GetDir() (dir string, err error) {
	return r.envParams.GetExeDir()
}
