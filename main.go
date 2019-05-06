package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	configFlag   = flag.String("config", "/opt/prometheus-actions/prometheus-actions.yaml", "Path to configuration file")
	logLevelFlag = flag.String("log-level", "debug", "Logging level")
)

func realMain() int {
	flag.Parse()
	log := logrus.New()
	if level, err := logrus.ParseLevel(*logLevelFlag); err == nil {
		log.SetLevel(level)
	}
	log.Formatter = &logrus.JSONFormatter{}
	log.Info("Starting prometheus-actions...")
	config, err := LoadConfig(*configFlag)
	if err != nil {
		log.Error(err)
		return 1
	}
	if err := config.Validate(); err != nil {
		log.Error(err)
		return 1
	}
	executor, err := NewExecutor(log, config)
	if err != nil {
		log.Error(err)
		return 1
	}
	if err := executor.Run(); err != nil {
		log.Error(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
