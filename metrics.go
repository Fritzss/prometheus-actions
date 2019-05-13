package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

var (
	cmdExecuteDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "prometheus_actions_command_execute_duration_seconds",
			Help: "The duration of the command execution in seconds.",
		},
		[]string{"action"},
	)
	cmdExecuteErrorsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "prometheus_actions_command_execute_errors_total",
			Help: "The number of command execution errors.",
		},
		[]string{"action"},
	)
	promRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "prometheus_actions_prometheus_request_duration_seconds",
			Help: "The duration of Prometheus request in seconds.",
		},
		[]string{"action"},
	)
	promRequestErrorsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "prometheus_actions_prometheus_request_errors_total",
			Help: "The number of Prometheus request errors.",
		},
		[]string{"action"},
	)
	buildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prometheus_actions_build_info",
			Help: "A metric with a constant '1' value labeled by version, revision, branch, and goversion from which Prometheus Actions was built.",
		},
		[]string{"version", "revision", "branch", "goversion"},
	)
)

func init() {
	prometheus.MustRegister(cmdExecuteDuration)
	prometheus.MustRegister(cmdExecuteErrorsCount)
	prometheus.MustRegister(promRequestErrorsCount)
	prometheus.MustRegister(promRequestDuration)
	prometheus.MustRegister(buildInfo)

	buildInfo.WithLabelValues(
		version.Version,
		version.Revision,
		version.Branch,
		version.GoVersion,
	).Set(1)
}
