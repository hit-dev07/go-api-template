package logging

import (
	"github.com/geometry-labs/go-service-template/config"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	//core.GetEnvironment()
	config.Vars.ConfigFile = "config.api.test"
	config.Vars.ConfigType = "yaml"
	config.Vars.ConfigPath = "../../envfiles"
	config.ConfigInit()
}

func TestInit(t *testing.T) {

	// Test file logging
	os.Setenv("LOG_LEVEL", "Info")
	os.Setenv("LOG_TO_FILE", "true")
	config.GetEnvironment()
	StartLoggingInit()
	log.Info("File log")

	// Test levels
	os.Setenv("LOG_LEVEL", "Panic")
	config.GetEnvironment()
	StartLoggingInit()
	log.Info("Should not log")

	os.Setenv("LOG_LEVEL", "FATAL")
	config.GetEnvironment()
	StartLoggingInit()
	log.Info("Should not log")

	os.Setenv("LOG_LEVEL", "ERROR")
	config.GetEnvironment()
	StartLoggingInit()
	log.Info("Should not log")

	os.Setenv("LOG_LEVEL", "WARN")
	config.GetEnvironment()
	StartLoggingInit()
	log.Warn("Warning")

	os.Setenv("LOG_LEVEL", "INFO")
	config.GetEnvironment()
	StartLoggingInit()
	log.Info("Info")

	os.Setenv("LOG_LEVEL", "DEBUG")
	config.GetEnvironment()
	StartLoggingInit()
	log.Debug("Debug")

	os.Setenv("LOG_LEVEL", "TRACE")
	config.GetEnvironment()
	StartLoggingInit()
	log.Trace("Trace")
}
