package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Executor struct {
	c      *Config
	f      *Fingerprint
	promQL v1.API
}

func NewExecutor(config *Config) (*Executor, error) {
	e := &Executor{
		c: config,
	}
	if err := e.setupFingerprint(); err != nil {
		return nil, err
	}
	if err := e.setupPromQL(); err != nil {
		return nil, err
	}
	if err := e.compileQueries(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Executor) compileQueries() error {
	for _, action := range e.c.Actions {
		ql, err := generateTemplate(action.Expr, action.String(), e.f)
		if err != nil {
			return err
		}
		action.compiledExpr = StandardizeSpaces(ql)
	}
	return nil
}

func (e *Executor) setupFingerprint() error {
	fingerprint, err := BuildFingerprint()
	if err != nil {
		return err
	}
	e.f = fingerprint
	return nil
}

func (e *Executor) setupPromQL() error {
	promCfg := api.Config{
		Address: e.c.PrometheusURL,
	}
	promCli, err := api.NewClient(promCfg)
	if err != nil {
		return err
	}
	q := v1.NewAPI(promCli)
	e.promQL = q
	return nil
}

func (e *Executor) ExecuteQuery(q string) (model.Value, error) {
	result, err := e.promQL.Query(context.Background(), q, time.Now())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Executor) MustExecuteCommand(result model.Value) (bool, error) {
	switch {
	case result.Type() == model.ValVector:
		samples := result.(model.Vector).Len()
		if samples > 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("Unexpected result type: %v", result.Type())
}

func (e *Executor) ExecuteCommand(command []string) error {
	var cmd *exec.Cmd
	ctx, cancel := context.WithTimeout(context.Background(), e.c.СommandTimeout)
	defer cancel()
	if len(command) == 1 {
		cmd = exec.CommandContext(ctx, command[0])
	} else {
		cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("Non-zero exit code: %v", err)
	}
	return nil
}

func (e *Executor) processActions() error {
	for _, action := range e.c.Actions {
		log.Printf("Querying '%s' for %s...", action.compiledExpr, action.String())
		result, err := e.ExecuteQuery(action.compiledExpr)
		if err != nil {
			log.Printf("Failed to query: %v", err)
			continue
		}
		ok, err := e.MustExecuteCommand(result)
		if err != nil {
			log.Printf("Failed to check query result: %v", err)
			continue
		}
		if !ok {
			continue
		}
		log.Printf("Executing '%s' for %s...", strings.Join(action.Command, " "), action.String())
		if err := e.ExecuteCommand(action.Command); err != nil {
			log.Printf("Failed to execute: %v", err)
			continue
		}
		log.Printf("Done with %s", action.String())
	}
	return nil
}

func (e *Executor) Run() error {
	for {
		err := e.processActions()
		if err != nil {
			return err
		}
		log.Printf("Sleeping for a %s...", e.c.QueryInterval)
		time.Sleep(e.c.QueryInterval)
	}
}
