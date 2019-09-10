package params

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("port", "8000")
	viper.SetDefault("logmode", "")
	viper.SetDefault("loglevel", "")
	viper.SetDefault("nodeid", "0")
	viper.BindEnv("port")
	viper.BindEnv("logmode")
	viper.BindEnv("loglevel")
	viper.BindEnv("nodeid")
}
