# Prometheus Actions

## Configuration example

```yaml
---
actionsInterval: 10m
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
