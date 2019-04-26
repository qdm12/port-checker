package params

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("listeningport", "8000")
	viper.SetDefault("logging", "human")
	viper.SetDefault("nodeid", "0")
	viper.BindEnv("listeningport")
	viper.BindEnv("logging")
	viper.BindEnv("nodeid")
}
