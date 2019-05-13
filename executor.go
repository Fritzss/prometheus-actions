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

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

var (
	// Simple protection to prevent Prometheus DDoS
	defaultRepeatDelay = time.Second
)

type Executor struct {
	mux        *http.ServeMux
	httpServer http.Server
	c          *Config
	f          *Fingerprint
	log        *logrus.Logger
	environ    []string
	promQL     v1.API
}

func NewExecutor(log *logrus.Logger, config *Config) (*Executor, error) {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    config.ListenAddress,
		Handler: mux,
	}
	e := &Executor{
		mux:        mux,
		httpServer: server,
		c:          config,
		log:        log,
		environ:    os.Environ(),
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
	e.f = fingerprint
	return err
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

func (e *Executor) processAction(action *Action) {
	logEntry := e.log.WithField("action", action.String())
	if limited := action.IsCooldownLimited(e.c.CooldownPeriod); limited {
		logEntry.Debug("Can't process due cooldown period")
		return
	}

	logEntry.Debugf("Querying '%s'...", action.compiledExpr)
	t0 := time.Now()
	result, err := e.ExecuteQuery(action.compiledExpr)
	promRequestDuration.WithLabelValues(action.Name).Observe(time.Since(t0).Seconds())
	if err != nil {
		promRequestErrorsCount.WithLabelValues(action.Name).Inc()
		logEntry.Errorf("Failed to query: %v", err)
		return
	}

	canExecute, err := e.CanExecuteCommand(result)
	if err != nil {
		logEntry.Errorf("Failed to check query result: %v", err)
		return
	}
	if !canExecute {
		return
	}

	logEntry.Infof("Executing '%s'...", strings.Join(action.Command, " "))
	action.lastExecTime = time.Now()

	t1 := time.Now()
	err = e.ExecuteCommand(action.Command)
	cmdExecuteDuration.WithLabelValues(action.Name).Observe(time.Since(t1).Seconds())
	if err != nil {
		cmdExecuteErrorsCount.WithLabelValues(action.Name).Inc()
		logEntry.Errorf("Failed to execute: %v", err)
		return
	}

	logEntry.Debug("Done")
}

func (e *Executor) processActions() {
	for _, action := range e.c.Actions {
		e.processAction(action)
		time.Sleep(defaultRepeatDelay)
	}
}

func (e *Executor) serveRequests() error {
	return e.httpServer.ListenAndServe()
}

func (e *Executor) registerHandlers() {
	e.mux.Handle("/metrics", prometheus.Handler())
}

func (e *Executor) Run(ctx context.Context) error {
	errCh := make(chan error)
	e.registerHandlers()
	go func() {
		errCh <- e.serveRequests()
	}()
	next := time.After(time.Second)
	for {
		select {
		case <-ctx.Done():
			e.httpServer.Shutdown(ctx)
			return nil
		case err := <-errCh:
			return err
		case <-next:
			e.processActions()
			next = time.After(e.c.RepeatInterval)
			e.log.Debugf("Sleeping for a %s...", e.c.RepeatInterval)
		}
	}
}
