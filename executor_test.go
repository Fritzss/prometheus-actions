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
