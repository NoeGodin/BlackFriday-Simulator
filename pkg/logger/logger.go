package logger

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger    *logrus.Logger
	debugMode bool
)

func init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func InitLogger() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug logging")
	flag.Parse()

	if debugMode {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.Debug("Debug mode enabled")
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}
