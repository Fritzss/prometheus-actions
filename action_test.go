package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Action(t *testing.T) {
	a := &Action{}
	assert.Error(t, a.Validate())

	a = &Action{
		Expr: "query",
	}
	assert.Error(t, a.Validate())

	a = &Action{
		Expr: "query",
		Command: []string{
			"cmd",
		},
	}
	assert.NoError(t, a.Validate())
}

func TestString(t *testing.T) {
	a := &Action{}
	assert.Contains(t, a.String(), "unnamed")

	a = &Action{
		Name: "name",
	}
	assert.Contains(t, a.String(), "name")
}

func TestIsCooldownLimited(t *testing.T) {
	a := &Action{}
	assert.False(t, a.IsCooldownLimited(10*time.Minute))

	a = &Action{
		lastExecTime: time.Now().Add(-15 * time.Minute),
	}
	assert.False(t, a.IsCooldownLimited(10*time.Minute))

	a = &Action{
		lastExecTime: time.Now(),
	}
	assert.True(t, a.IsCooldownLimited(10*time.Minute))
}
