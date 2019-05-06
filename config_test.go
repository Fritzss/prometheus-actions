package main

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("not-found")
	if err == nil {
		t.Error("Must be an error, but got a nil")
	}
	_, err = LoadConfig("test_data/config_invalid.yaml")
	if err == nil {
		t.Error("Must be an error, but got a nil")
	}
	_, err = LoadConfig("test_data/config_valid.yaml")
	if err != nil {
		t.Error(err)
	}
}

func TestValidate_Config(t *testing.T) {
	actions := []*Action{
		&Action{
			Name: "name",
			Command: []string{
				"cmd",
			},
		},
	}
	c := &Config{}
	if err := c.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	c = &Config{
		Actions: actions,
	}
	if err := c.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	c = &Config{
		Actions:         actions,
		ActionsInterval: time.Minute,
	}
	if err := c.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	c = &Config{
		Actions:         actions,
		ActionsInterval: time.Minute,
		СommandTimeout:  time.Minute,
	}
	if err := c.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	c.Actions[0].Expr = "query"
	if err := c.Validate(); err != nil {
		t.Error(err)
	}
}
