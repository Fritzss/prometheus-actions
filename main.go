package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

var (
	configFlag   = flag.String("config", "/opt/prometheus-actions/prometheus-actions.yaml", "Path to configuration file")
	logLevelFlag = flag.String("log-level", "debug", "Logging level")
	versionFlag  = flag.Bool("version", false, "Prints version and exit")
)

func realMain() int {
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Print("prometheus-actions"))
		return 0
	}

	log := logrus.New()
	if level, err := logrus.ParseLevel(*logLevelFlag); err == nil {
		log.SetLevel(level)
	}
	log.Formatter = &logrus.JSONFormatter{}
	log.Infof("Starting prometheus-actions %s...", version.Version)

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
	if err := executor.Run(context.Background()); err != nil {
		log.Error(err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
