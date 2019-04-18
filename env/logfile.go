package env

import (
	"io"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogFileConfig keeps configuration settings required for lumberjack
type LogFileConfig struct {
	Filename  string
	MaxSize   int
	Level     logrus.Level
	Formatter logrus.Formatter
}

// LogFileHook structure implements logfile hook for logrus
type LogFileHook struct {
	Config    LogFileConfig
	logWriter io.Writer
}

// NewLogFileHook instantiates hook and implements Hook interface
func NewLogFileHook(config LogFileConfig) (logrus.Hook, error) {
	hook := LogFileHook{
		Config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename: config.Filename,
		MaxSize:  config.MaxSize,
	}

	return &hook, nil
}

// Levels defines and returns the log levels under which logrus fires the "Fire"
func (hook *LogFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

// Fire implements actual write()
func (hook *LogFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.logWriter.Write(b)

	return nil
}
