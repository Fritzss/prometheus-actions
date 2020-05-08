package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	errorResult  = `{`
	matrixResult = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	emptyResult  = `{"status":"success","data":{"resultType":"vector","result":[]}}`
	fullResult   = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","instance":"127.0.0.1:9100","job":"test"},"value":[1557382679.814,"1"]}]}}`
)

var (
	mock = &promMock{}
)

type promMock struct {
	result string
}

func (p *promMock) start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(p.result))
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})
	ch := make(chan error)
	go func() {
		ch <- http.ListenAndServe("127.0.0.1:9001", nil)
	}()
	select {
	case err := <-ch:
		return err
	case <-time.NewTicker(time.Second).C:
	}
	return nil
}

func TestExecuteCommand(t *testing.T) {
	e := &Executor{
		log: logrus.New(),
		c: &Config{
			CommandTimeout: 100 * time.Millisecond,
		},
	}
	if err := e.ExecuteCommand([]string{"whoami"}); err != nil {
		t.Error(err)
	}
	if err := e.ExecuteCommand([]string{"sleep", "0.5"}); err == nil {
		t.Error("Must be timeout")
	}
	if err := e.ExecuteCommand([]string{"exit", "1"}); err == nil {
		t.Error("Must be an error")
	}
}

func TestNewExecutor(t *testing.T) {
	config, err := LoadConfig("fixtures/config_valid.yaml")
	if err != nil {
		t.Fatal(err)
	}
	log := logrus.New()
	config.Actions[0].Expr = "{{ ."
	_, err = NewExecutor(log, config)
	if err == nil {
		t.Error("Must be an error")
	}
	config.Actions[0].Expr = "up"
	config.PrometheusURL = "@#$%^&*()"
	_, err = NewExecutor(log, config)
	if err == nil {
		t.Error("Must be an error")
	}
}

func TestRun(t *testing.T) {
	err := mock.start()
	if err != nil {
		t.Fatal(err)
	}

	testRun(t, ":9301", fullResult)
	testRun(t, ":9302", emptyResult)
	testRun(t, ":9303", errorResult)
	testRun(t, ":9304", matrixResult)
}

func testRun(t *testing.T, listenAddress, result string) {
	mock.result = result
	log := logrus.New()
	ctx, cancel := context.WithCancel(context.Background())

	defaultRepeatDelay = time.Millisecond
	config, err := LoadConfig("fixtures/config_valid.yaml")
	if err != nil {
		t.Fatal(err)
	}
	config.ListenAddress = listenAddress

	executor, err := NewExecutor(log, config)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan error)
	go func() {
		ch <- executor.Run(ctx)
	}()

	select {
	case err := <-ch:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.NewTicker(time.Second).C:
		err := executor.serveRequests()
		if err == nil {
			t.Error("Must be an error")
		}
	}

	executor.processActions()
	executor.c.CooldownPeriod = 5 * time.Minute
	executor.processActions()
	cancel()
}
