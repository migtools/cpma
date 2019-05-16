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

	expectedHomeDir, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name                   string
		expectedHomeDir        string
		expectedConfigFilePath string
	}{
		{
			name:                   "Init config",
			expectedHomeDir:        expectedHomeDir,
			expectedConfigFilePath: ConfigFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualHomeDir := viperConfig.GetString("home")
			assert.Equal(t, tc.expectedHomeDir, actualHomeDir)
			actualConfigFilePath := viperConfig.ConfigFileUsed()
			assert.Equal(t, tc.expectedConfigFilePath, actualConfigFilePath)
		})
	}
}

func TestInitLogger(t *testing.T) {
	expectedFileHook, err := NewLogFileHook(
		LogFileConfig{
			Filename: logFile,
			MaxSize:  5, // MiB
			Level:    logrus.InfoLevel,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint: true,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name              string
		expectedLogLevel  logrus.Level
		expectedFormatter *logrus.TextFormatter
		expectedFileHook  logrus.Hook
		debugLevel        bool
	}{
		{
			name:             "init logger",
			expectedLogLevel: logrus.InfoLevel,
			expectedFormatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC822,
			},
			expectedFileHook: expectedFileHook,
			debugLevel:       false,
		},
		{
			name:             "init logger",
			expectedLogLevel: logrus.DebugLevel,
			expectedFormatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC822,
			},
			expectedFileHook: expectedFileHook,
			debugLevel:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			viperConfig.Set("debug", tc.debugLevel)
			InitLogger()
			logger := logrus.StandardLogger()
			if tc.debugLevel {
				assert.Equal(t, tc.expectedLogLevel, logger.GetLevel())
			} else {
				assert.Equal(t, tc.expectedLogLevel, logger.GetLevel())
				assert.Equal(t, tc.expectedFormatter, logger.Formatter)
				assert.Equal(t, tc.expectedFileHook, logger.Hooks[logrus.InfoLevel][0])
			}
		})
	}
}
