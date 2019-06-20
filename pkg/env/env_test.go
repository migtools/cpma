package env

import (
	"os"
	"testing"
	"time"

	"github.com/fusor/cpma/pkg/api"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
)

func TestInitConfig(t *testing.T) {
	ConfigFile = "testdata/cpma-config.yml"
	api.Client = &kubernetes.Clientset{}
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

	consoleFormatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
		ForceColors:     true,
	}

	expectedStderrHook := &ConsoleWriterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		},
		Formatter: consoleFormatter,
	}

	expectedStdoutHook := &ConsoleWriterHook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.WarnLevel,
		},
		Formatter: consoleFormatter,
	}

	testCases := []struct {
		name               string
		expectedLogLevel   logrus.Level
		expectedFormatter  *logrus.TextFormatter
		expectedFileHook   logrus.Hook
		expectedStderrHook logrus.Hook
		expectedStdoutHook logrus.Hook
		debugLevel         bool
	}{
		{
			name:             "init logger",
			expectedLogLevel: logrus.InfoLevel,
			expectedFormatter: &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC822,
			},
			expectedFileHook:   expectedFileHook,
			expectedStderrHook: expectedStderrHook,
			expectedStdoutHook: expectedStdoutHook,
			debugLevel:         false,
		},
		{
			name:             "init logger and set log level to debug",
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
			viperConfig.Set("consolelogs", true)
			InitLogger()
			logger := logrus.StandardLogger()
			if tc.debugLevel {
				assert.Equal(t, tc.expectedLogLevel, logrus.GetLevel())
			} else {
				assert.Equal(t, tc.expectedLogLevel, logger.GetLevel())

				assert.Equal(t, tc.expectedFileHook, logger.Hooks[logrus.InfoLevel][0])
				assert.Equal(t, tc.expectedStdoutHook, logger.Hooks[logrus.InfoLevel][1])

				assert.Equal(t, tc.expectedFileHook, logger.Hooks[logrus.ErrorLevel][0])
				assert.Equal(t, tc.expectedStderrHook, logger.Hooks[logrus.ErrorLevel][1])
			}
		})
	}
}
