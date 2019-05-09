package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultListenAddress = "0.0.0.0:9333"
)

type Config struct {
	ListenAddress   string        `yaml:"listenAddress"`
	PrometheusURL   string        `yaml:"prometheusURL"`
	ActionsInterval time.Duration `yaml:"actionsInterval"`
	СommandTimeout  time.Duration `yaml:"commandTimeout"`
	CooldownPeriod  time.Duration `yaml:"cooldownPeriod"`
	Actions         []*Action
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
	config.SpecifyDefaults()
	return config, nil
}

func (c *Config) SpecifyDefaults() {
	if len(c.ListenAddress) == 0 {
		c.ListenAddress = defaultListenAddress
	}
}

func (c *Config) Validate() error {
	if len(c.Actions) == 0 {
		return errors.New("Actions must be specified")
	}
	if c.ActionsInterval <= time.Second {
		return fmt.Errorf("ActionsInterval must be greater than second")
	}
	if c.СommandTimeout <= time.Second {
		return fmt.Errorf("СommandTimeout must be greater than second")
	}
	uniqueActions := make(map[string]bool)
	for i, action := range c.Actions {
		err := action.Validate()
		if err != nil {
			return fmt.Errorf("Action %d error: %v", i, err)
		}
		if _, ok := uniqueActions[action.Name]; ok {
			return fmt.Errorf("Duplicate of %s action", action.Name)
		}
		uniqueActions[action.Name] = true
	}
	return nil
}
