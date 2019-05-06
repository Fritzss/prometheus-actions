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
	CooldownPeriod time.Duration `yaml:"cooldownPeriod"`
	Actions        []*Action
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
