---
title: "Configuration"
linkTitle: "Configuration"
weight: 2
description: >
  Learn how to configure Watchdog.
---

Watchdog configuration uses the [YAML](https://yaml.org/) format.

The file to be edited can be found in:
1. `/etc/watchdog/config.yaml` on \*nix systems when Watchdog is executed as root
2. `~/.watchdog/config.yaml` on \*nix systems when Watchdog is executed as non-root
2. `./config.yaml` on other systems

Configuration example:

```yaml
concurrent: 4
logs_level: "warning"
```