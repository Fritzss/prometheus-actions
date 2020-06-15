#!/usr/bin/env bash
# ref: https://hub.docker.com/r/prom/prometheus/tags

set -euo pipefail

echo "Starting Prometheus..."
docker rm -f dev-prom &> /dev/null || :
docker run -d --name=dev-prom -p 9090:9090 prom/prometheus:v2.19.0 &> /dev/null
echo "...done. Run \"docker rm -f dev-prom\" to clean up the container."
echo
