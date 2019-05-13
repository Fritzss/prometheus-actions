# Prometheus Actions

[![BuildStatus Widget]][BuildStatus Result]
[![codecov](https://codecov.io/gh/leominov/prometheus-actions/branch/master/graph/badge.svg)](https://codecov.io/gh/leominov/prometheus-actions)

[BuildStatus Result]: https://travis-ci.com/leominov/prometheus-actions
[BuildStatus Widget]: https://travis-ci.com/leominov/prometheus-actions.svg?branch=master

## Configuration example

```yaml
---
repeatInterval: 10m
cooldownPeriod: 30m
commandTimeout: 5m
prometheusURL: http://prometheus.local/
listenAddress: 0.0.0.0:9333
actions:
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

## Links

* [Quering Prometheus](https://prometheus.io/docs/prometheus/latest/querying/basics/)
