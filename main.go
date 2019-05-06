package main

import (
	"flag"
	"log"
	"os"
)

var (
	configFlag = flag.String(
		"config",
		"/opt/prometheus-actions/prometheus-actions.yaml",
		"Path to configuration file",
	)
)

func realMain() int {
	flag.Parse()
	log.Println("Starting prometheus-actions...")
	config, err := LoadConfig(*configFlag)
	if err != nil {
		log.Println(err)
		return 1
	}
	if err := config.Validate(); err != nil {
		log.Println(err)
		return 1
	}
	executor, err := NewExecutor(config)
	if err != nil {
		log.Println(err)
		return 1
	}
	if err := executor.Run(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
