package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("not-found")
	assert.Error(t, err)

	_, err = LoadConfig("fixtures/config_invalid.yaml")
	assert.Error(t, err)

	_, err = LoadConfig("fixtures/config_valid.yaml")
	assert.NoError(t, err)
}

func TestSpecifyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.SpecifyDefaults()
	assert.Equal(t, defaultListenAddress, cfg.ListenAddress)
}

func TestValidate_Config(t *testing.T) {
	actions := []*Action{
		{
			Name: "name",
			Command: []string{
				"cmd",
			},
		},
	}
	c := &Config{}
	assert.Error(t, c.Validate())

	c = &Config{
		Actions: actions,
	}
	assert.Error(t, c.Validate())

	c = &Config{
		Actions:        actions,
		RepeatInterval: time.Minute,
	}
	assert.Error(t, c.Validate())

	c = &Config{
		Actions:        actions,
		RepeatInterval: time.Minute,
		CommandTimeout: time.Minute,
	}
	assert.Error(t, c.Validate())

	c.Actions[0].Expr = "query"
	assert.NoError(t, c.Validate())

	c.Actions = append(c.Actions, c.Actions[0])
	assert.Error(t, c.Validate())
}
