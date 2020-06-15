package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, e.ExecuteCommand([]string{"whoami"}, nil))
	assert.Error(t, e.ExecuteCommand([]string{"sleep", "0.5"}, nil))
	assert.Error(t, e.ExecuteCommand([]string{"exit", "1"}, nil))
}

func TestNewExecutor(t *testing.T) {
	config, err := LoadConfig("fixtures/config_valid.yaml")
	assert.NoError(t, err)

	log := logrus.New()
	config.Actions[0].Expr = "{{ ."
	_, err = NewExecutor(log, config)
	assert.Error(t, err)

	config.Actions[0].Expr = "up"
	config.PrometheusURL = "@#$%^&*()"
	_, err = NewExecutor(log, config)
	assert.Error(t, err)
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
		assert.Error(t, executor.serveRequests())
	}

	executor.processActions()
	executor.c.CooldownPeriod = 5 * time.Minute
	executor.processActions()
	cancel()
}

func TestExecutor_ParseQueryResult(t *testing.T) {
	prometheusURL := "http://127.0.0.1:9090"
	_, err := http.Get(fmt.Sprintf("%s/-/healthy", prometheusURL))
	if err != nil {
		t.Skip("Run Prometheus before tests")
	}

	config := &Config{
		PrometheusURL: prometheusURL,
	}

	ex, err := NewExecutor(logrus.New(), config)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ex.ExecuteQuery(`count(up) > 0`)
	assert.NoError(t, err)

	labelSetSlice, ok, err := ex.ParseQueryResult(result)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 0, len(labelSetSlice))

	result, err = ex.ExecuteQuery(`up`)
	assert.NoError(t, err)

	labelSetSlice, ok, err = ex.ParseQueryResult(result)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 1, len(labelSetSlice))

	for _, labelSet := range labelSetSlice {
		for labelKey, labelValue := range labelSet {
			if labelKey == ("__name__") {
				assert.Equal(t, model.LabelValue("up"), labelValue)
			}
		}
	}

	result, err = ex.ExecuteQuery(`count(up) > 10`)
	assert.NoError(t, err)

	labelSetSlice, ok, err = ex.ParseQueryResult(result)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 0, len(labelSetSlice))

	result, err = ex.ExecuteQuery(`count(up2) > 0`)
	assert.NoError(t, err)

	labelSetSlice, ok, err = ex.ParseQueryResult(result)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, 0, len(labelSetSlice))
}
