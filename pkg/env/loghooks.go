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

// ConsoleWriterHook is a hook that writes logs of specified LogLevels to specified Writer
type ConsoleWriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
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

	_, err = hook.logWriter.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *ConsoleWriterHook) Fire(entry *logrus.Entry) error {
	b, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.Writer.Write(b)
	if err != nil {
		return err
	}

	return err
}

// Levels define on which log levels this hook would trigger
func (hook *ConsoleWriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}
