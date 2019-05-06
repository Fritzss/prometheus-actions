# Prometheus Actions

## Configuration example

```yaml
---
queryInterval: 10m
cooldownPeriod: 30m
commandTimeout: 5m
prometheusURL: http://prometheus.local/
actions:
  - name: Docker GC
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
