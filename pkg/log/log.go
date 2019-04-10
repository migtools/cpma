package log

import (
	"io"
	"os"

	"github.com/fusor/cpma/pkg/config"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

// Logger defines a set of methods for writing application logs. Derived from and
// inspired by logrus.Entry.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// DefaultLogger default
var DefaultLogger *logrus.Logger

func init() {
	DefaultLogger = newLogrusLogger(config.Config())
}

func newLogrusLogger(cfg *viper.Viper) *logrus.Logger {
	l := logrus.New()

	f, err := os.OpenFile("cpma.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		l.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)

	l.Out = mw

	l.Level = logrus.InfoLevel

	l.Println("CPMA Log started")

	return l
}

// SetDebugLogLevel set loglevel to debug
func SetDebugLogLevel() {
	DefaultLogger.Level = logrus.DebugLevel
}

type Fields map[string]interface{}

func (f Fields) With(k string, v interface{}) Fields {
	f[k] = v
	return f
}

func (f Fields) WithFields(f2 Fields) Fields {
	for k, v := range f2 {
		f[k] = v
	}
	return f
}

func WithFields(fields Fields) Logger {
	return DefaultLogger.WithFields(logrus.Fields(fields))
}

// Debug package-level convenience method.
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

// Debugf package-level convenience method.
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Debugln package-level convenience method.
func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}

// Error package-level convenience method.
func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

// Errorf package-level convenience method.
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Errorln package-level convenience method.
func Errorln(args ...interface{}) {
	DefaultLogger.Errorln(args...)
}

// Fatal package-level convenience method.
func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

// Fatalf package-level convenience method.
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// Fatalln package-level convenience method.
func Fatalln(args ...interface{}) {
	DefaultLogger.Fatalln(args...)
}

// Panic package-level convenience method.
func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

// Panicf package-level convenience method.
func Panicf(format string, args ...interface{}) {
	DefaultLogger.Panicf(format, args...)
}

// Panicln package-level convenience method.
func Panicln(args ...interface{}) {
	DefaultLogger.Panicln(args...)
}

// Print package-level convenience method.
func Print(args ...interface{}) {
	DefaultLogger.Print(args...)
}

// Printf package-level convenience method.
func Printf(format string, args ...interface{}) {
	DefaultLogger.Printf(format, args...)
}

// Println package-level convenience method.
func Println(args ...interface{}) {
	DefaultLogger.Println(args...)
}
