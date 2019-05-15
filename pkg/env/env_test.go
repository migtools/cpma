package env

import (
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	ConfigFile = "testdata/test-cpma-config.yml"
	err := InitConfig()
	if err != nil {
		t.Fatal(err)
	}

	expectedHomeDir, _ := homedir.Dir()
	actualHomeDir := viperConfig.GetString("home")
	assert.Equal(t, expectedHomeDir, actualHomeDir)

	expectedConfigFilePath := ConfigFile
	actualConfigFilePath := viperConfig.ConfigFileUsed()
	assert.Equal(t, expectedConfigFilePath, actualConfigFilePath)
}

func TestInitLogger(t *testing.T) {
	InitLogger()
	logger := logrus.StandardLogger()

	// Test if info level is set by default
	expectedLogLevel := logrus.InfoLevel
	assert.Equal(t, expectedLogLevel, logger.GetLevel())

	// Test formatter
	expectedFormatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	}
	assert.Equal(t, expectedFormatter, logger.Formatter)

	// Test if hook is set right
	expectedFileHook, _ := NewLogFileHook(
		LogFileConfig{
			Filename: logFile,
			MaxSize:  5, // MiB
			Level:    logrus.InfoLevel,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint: true,
			},
		},
	)
	assert.Equal(t, expectedFileHook, logger.Hooks[logrus.InfoLevel][0])

	// Test if debug is set from config
	viperConfig.Set("debug", true)
	InitLogger()
	expectedLogLevel = logrus.DebugLevel
	assert.Equal(t, expectedLogLevel, logger.GetLevel())
}
