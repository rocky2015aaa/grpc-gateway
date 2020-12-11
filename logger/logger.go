// Package logger works for logging system
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

const (
	FUNC_NAME = ""
	SEVERITY  = ""

	PANIC_LEVEL = logrus.PanicLevel
	FATAL_LEVEL = logrus.FatalLevel
	ERROR_LEVEL = logrus.ErrorLevel
	WARN_LEVEL  = logrus.WarnLevel
	INFO_LEVEL  = logrus.InfoLevel
	DEBUG_LEVEL = logrus.DebugLevel
	TRACE_LEVEL = logrus.TraceLevel
)

func init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(INFO_LEVEL)
}
