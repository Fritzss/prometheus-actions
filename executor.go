package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Executor struct {
	c       *Config
	f       *Fingerprint
	log     *logrus.Logger
	environ []string
	promQL  v1.API
}

func NewExecutor(log *logrus.Logger, config *Config) (*Executor, error) {
	e := &Executor{
		c:       config,
		log:     log,
		environ: os.Environ(),
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

func (e *Executor) CanExecuteCommand(result model.Value) (bool, error) {
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
	cmd.Env = e.environ
	cmd.Stderr = e.log.WithField("src", "cmd").WriterLevel(logrus.ErrorLevel)
	cmd.Stdout = e.log.WithField("src", "cmd").WriterLevel(logrus.DebugLevel)
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("Non-zero exit code: %v", err)
	}
	return nil
}

func (e *Executor) processActions() {
	for _, action := range e.c.Actions {
		logEntry := e.log.WithField("action", action.String())
		if limited := action.IsCooldownLimited(e.c.CooldownPeriod); limited {
			logEntry.Infof("Can't process due cooldown period")
			continue
		}
		logEntry.Debugf("Querying '%s'...", action.compiledExpr)
		result, err := e.ExecuteQuery(action.compiledExpr)
		if err != nil {
			logEntry.Errorf("Failed to query: %v", err)
			continue
		}
		canExecute, err := e.CanExecuteCommand(result)
		if err != nil {
			logEntry.Errorf("Failed to check query result: %v", err)
			continue
		}
		if !canExecute {
			continue
		}
		logEntry.Debugf("Executing '%s'...", strings.Join(action.Command, " "))
		action.lastExecTime = time.Now()
		if err := e.ExecuteCommand(action.Command); err != nil {
			logEntry.Errorf("Failed to execute: %v", err)
			continue
		}
		logEntry.Debug("Done")
	}
}

func (e *Executor) serveRequests() error {
	http.Handle("/metrics", prometheus.Handler())
	return http.ListenAndServe(e.c.ListenAddress, nil)
}

func (e *Executor) Run() error {
	errCh := make(chan error)
	go func() {
		errCh <- e.serveRequests()
	}()
	next := time.After(time.Second)
	for {
		select {
		case err := <-errCh:
			return err
		case <-next:
			e.processActions()
			next = time.After(e.c.QueryInterval)
			e.log.Debugf("Sleeping for a %s...", e.c.QueryInterval)
		}
	}
}
