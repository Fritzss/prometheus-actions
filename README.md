# Prometheus Actions

[![Build Status](https://travis-ci.com/leominov/prometheus-actions.svg?token=tyxzVzn67Z9UV2wuxhSV&branch=master)](https://travis-ci.com/leominov/prometheus-actions)
[![codecov](https://codecov.io/gh/leominov/prometheus-actions/branch/master/graph/badge.svg)](https://codecov.io/gh/leominov/prometheus-actions)

## Configuration example

```yaml
---
repeatInterval: 10m
cooldownPeriod: 30m
commandTimeout: 5m
prometheusURL: http://prometheus.local/
listenAddress: 0.0.0.0:9333
actions:
  # Example 1: Doocker GC based on free space
  - name: Docker GC
    # Only Vectors supported for now
    expr: |
      (
        node_filesystem_free{instance="{{ .Hostname }}", mountpoint="/var/lib/docker"} /
        node_filesystem_size{instance="{{ .Hostname }}", mountpoint="/var/lib/docker"}
      ) * 100 < 10
    command:
      - bash
      - -c
      - "FORCE_IMAGE_REMOVAL=1 GRACE_PERIOD_SECONDS=3600 /usr/sbin/docker-gc"
  # Example 2: Restart GitLab runner before alert was fired
  - name: GitLab Runner Self-healing
    expr: |
      ALERTS{instance="{{ .Hostname }}", alertname="GitLabRunnerDown", alertstate="pending"} == 1
    command:
      - systemctl
      - restart
      - gitlab-runner
  # Example 3: Runs gitlabsos when GitLab server goes down
  # ref: https://gitlab.com/gitlab-com/support/toolbox/gitlabsos
  - name: GitLab SOS
    expr: |
      ALERTS{instance="{{ .Hostname }}", alertname="GitLabServerDown"} == 1
    command:
      - gitlabsos
      - --dir
      - /opt/gitlabsos
```

## Template variables

* `OSName` – One of ubuntu, linuxmint, and so on;
* `OSFamily` – One of debian, rhel, and so on;
* `OSVersion` – Version of the complete OS;
* `KernelName` – One of darwin, linux, and so on;
* `KernelVersion` – Version of the OS kernel (if available);
* `Hostname`.

## Template functions

Replace, default, length, lower, upper, urlencode, trim, yesno.

See [template_test.go](template_test.go) for examples.

## Metrics

* `prometheus_actions_build_info`
* `prometheus_actions_command_execute_duration_seconds`
* `prometheus_actions_command_execute_errors_total`
* `prometheus_actions_prometheus_request_duration_seconds`
* `prometheus_actions_prometheus_request_errors_total`

## Links

* [Quering Prometheus](https://prometheus.io/docs/prometheus/latest/querying/basics/)
