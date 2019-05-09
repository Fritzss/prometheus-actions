package main

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestExecuteCommand(t *testing.T) {
	e := &Executor{
		log: logrus.New(),
		c: &Config{
			СommandTimeout: 100 * time.Millisecond,
		},
	}
	if err := e.ExecuteCommand([]string{"whoami"}); err != nil {
		t.Error(err)
	}
	if err := e.ExecuteCommand([]string{"sleep", "0.5"}); err == nil {
		t.Error("Must be timeout")
	}
	if err := e.ExecuteCommand([]string{"exit", "1"}); err == nil {
		t.Error("Must be an error")
	}
}

func TestNewExecutor(t *testing.T) {
	config, err := LoadConfig("fixtures/config_valid.yaml")
	if err != nil {
		t.Fatal(err)
	}
	log := logrus.New()
	config.Actions[0].Expr = "{{ ."
	_, err = NewExecutor(log, config)
	if err == nil {
		t.Error("Must be an error")
	}
	config.Actions[0].Expr = "up"
	config.PrometheusURL = "@#$%^&*()"
	_, err = NewExecutor(log, config)
	if err == nil {
		t.Error("Must be an error")
	}
}

func TestRun(t *testing.T) {
	config, err := LoadConfig("fixtures/config_valid.yaml")
	if err != nil {
		t.Fatal(err)
	}
	log := logrus.New()
	executor, err := NewExecutor(log, config)
	if err != nil {
		t.Fatal(err)
	}
	defaultRepeatDelay = time.Millisecond
	ch := make(chan error)
	go func() {
		ch <- executor.Run()
	}()
	select {
	case err := <-ch:
		t.Fatal(err)
	case <-time.NewTicker(time.Second).C:
		err := executor.serveRequests()
		if err == nil {
			t.Error("Must be an error")
		}
	}
}
