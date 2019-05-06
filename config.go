package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	PrometheusURL  string        `yaml:"prometheusURL"`
	QueryInterval  time.Duration `yaml:"queryInterval"`
	СommandTimeout time.Duration `yaml:"commandTimeout"`
	Actions        []*Action
}

type Action struct {
	Name         string
	Command      []string
	Expr         string
	compiledExpr string
}

func LoadConfig(filename string) (*Config, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(out, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (a *Action) Validate() error {
	if len(a.Expr) == 0 {
		return errors.New("Action.Expr must be specified")
	}
	if len(a.Command) == 0 {
		return errors.New("Action.Command must be specified")
	}
	return nil
}

func (c *Config) Validate() error {
	if len(c.Actions) == 0 {
		return errors.New("Actions must be specified")
	}
	if c.QueryInterval <= time.Second {
		return fmt.Errorf("QueryInterval must be greater than second")
	}
	if c.СommandTimeout <= time.Second {
		return fmt.Errorf("СommandTimeout must be greater than second")
	}
	for i, action := range c.Actions {
		err := action.Validate()
		if err != nil {
			return fmt.Errorf("Action %d error: %v", i, err)
		}
	}
	return nil
}

func (a *Action) String() string {
	if len(a.Name) == 0 {
		return "unnamed"
	}
	return a.Name
}
