---
title: "Promtail"
linkTitle: "Promtail"
weight: 5
description: >
  Configuring Promtail.
---

Promtail is configured in a YAML file (usually referred to as config.yaml) which contains information on the Promtail server, where positions are stored, and how to scrape logs from files.

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://localhost:3100/loki/api/v1/push

scrape_configs:
- job_name: watchdog
  static_configs:
  - targets:
      - localhost
    labels:
      job: watchdog
      __path__: /var/log/watchdog/watchdog.log
```