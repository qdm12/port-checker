package params

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"port-checker/pkg/logging"

	"github.com/spf13/viper"
)

// GetListeningPort obtains and checks the listening port from Viper (env variable or config file, etc.)
func GetListeningPort() (listeningPort string) {
	listeningPort = viper.GetString("listeningPort")
	value, err := strconv.Atoi(listeningPort)
	if err != nil {
		logging.Fatal("listening port %s is not a valid integer", listeningPort)
	} else if value < 1 {
		logging.Fatal("listening port %s cannot be lower than 1", listeningPort)
	} else if value < 1024 {
		if os.Geteuid() == 0 {
			logging.Warn("listening port %s allowed to be in the reserved system ports range as you are running as root", listeningPort)
		} else if os.Geteuid() == -1 {
			logging.Warn("listening port %s allowed to be in the reserved system ports range as you are running in Windows", listeningPort)
		} else {
			logging.Fatal("listening port %s cannot be in the reserved system ports range (1 to 1023) when running without root", listeningPort)
		}
	} else if value > 65535 {
		logging.Fatal("listening port %s cannot be higher than 65535", listeningPort)
	} else if value > 49151 {
		// dynamic and/or private ports.
		logging.Warn("listening port %s is in the dynamic/private ports range (above 49151)", listeningPort)
	}
	return listeningPort
}

// GetRootURL obtains and checks the root URL from Viper (env variable or config file, etc.)
func GetRootURL() string {
	rootURL := viper.GetString("rooturl")
	if strings.ContainsAny(rootURL, " .?~#") {
		logging.Fatal("root URL %s contains invalid characters", rootURL)
	}
	rootURL = strings.ReplaceAll(rootURL, "//", "/")
	return strings.TrimSuffix(rootURL, "/") // already have / from paths of router
}

// GetDir obtains the executable directory
func GetDir() (dir string) {
	ex, err := os.Executable()
	if err != nil {
		logging.Fatal("%s", err)
	}
	return filepath.Dir(ex)
}

// GetLoggerMode obtains the logging mode from Viper (env variable or config file, etc.)
func GetLoggerMode() logging.Mode {
	s := viper.GetString("logging")
	return logging.ParseMode(s)
}

// GetLoggerLevel obtains the logging level from Viper (env variable or config file, etc.)
func GetLoggerLevel() logging.Level {
	s := viper.GetString("loglevel")
	return logging.ParseLevel(s)
}

// GetNodeID obtains the node instance ID from Viper (env variable or config file, etc.)
func GetNodeID() int {
	nodeID := viper.GetString("nodeid")
	value, err := strconv.Atoi(nodeID)
	if err != nil {
		logging.Fatal("Node ID %s is not a valid integer", nodeID)
	}
	return value
}
