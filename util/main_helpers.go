package util

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"gitlab.com/Fratbe/addglee/src/server/logger"
)

const (
	FN_WATCHING_SIGNAL = "WatchingSignal"
)

func SetConfig(conf *string) {
	flag.StringVar(conf, "", "", "")
	flag.Parse()
	if *conf == "" || !strings.Contains(*conf, "") {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

// watchingSignal monitors two syscall signal(SIGUSR1, SIGUSR2)
// to manage log level
func WatchingSignal() {
	sigusr1 := make(chan os.Signal, 1)
	sigusr2 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)
	signal.Notify(sigusr2, syscall.SIGUSR2)
	for {
		select {
		case <-sigusr1:
			nextLogLevel := logger.Log.GetLevel() + 1
			if nextLogLevel <= logger.TRACE_LEVEL {
				logger.Log.SetLevel(nextLogLevel)
			}
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_WATCHING_SIGNAL, logger.SEVERITY: logger.INFO_LEVEL}).Infoln("log level increased:", logger.Log.GetLevel())
		case <-sigusr2:
			nextLogLevel := logger.Log.GetLevel() - 1
			if nextLogLevel >= logger.PANIC_LEVEL {
				logger.Log.SetLevel(nextLogLevel)
			}
			logger.Log.WithFields(logrus.Fields{logger.FUNC_NAME: FN_WATCHING_SIGNAL, logger.SEVERITY: logger.INFO_LEVEL}).Infoln("log level decreased:", logger.Log.GetLevel())
		}
	}
}
