package main

import (
	"testing"
	"time"
)

func TestValidate_Action(t *testing.T) {
	a := &Action{}
	if err := a.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	a = &Action{
		Expr: "query",
	}
	if err := a.Validate(); err == nil {
		t.Error("Must be an error, but got nil")
	}
	a = &Action{
		Expr: "query",
		Command: []string{
			"cmd",
		},
	}
	if err := a.Validate(); err != nil {
		t.Error(err)
	}
}

func TestString(t *testing.T) {
	a := &Action{}
	if a.String() != "unnamed" {
		t.Errorf("Must be unnamed, but got %s", a.String())
	}
	a = &Action{
		Name: "name",
	}
	if a.String() != "name" {
		t.Errorf("Must be name, but got %s", a.String())
	}
}

func TestIsCooldownLimited(t *testing.T) {
	a := &Action{}
	if limited := a.IsCooldownLimited(10 * time.Minute); limited {
		t.Error("Must be false, but got true")
	}
	a = &Action{
		lastExecTime: time.Now().Add(-15 * time.Minute),
	}
	if limited := a.IsCooldownLimited(10 * time.Minute); limited {
		t.Error("Must be false, but got true")
	}
	a = &Action{
		lastExecTime: time.Now(),
	}
	if limited := a.IsCooldownLimited(10 * time.Minute); !limited {
		t.Error("Must be true, but got false")
	}
}
