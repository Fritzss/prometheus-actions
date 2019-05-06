package main

import (
	"errors"
	"time"
)

type Action struct {
	Name         string
	Command      []string
	Expr         string
	compiledExpr string
	lastExecTime time.Time
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

func (a *Action) String() string {
	if len(a.Name) == 0 {
		return "unnamed"
	}
	return a.Name
}
