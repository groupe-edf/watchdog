---
title: "Logrotate"
linkTitle: "Logrotate"
weight: 4
description: >
  Release binaries and docker image.
---

logrotate is designed to ease administration of systems that generate large numbers of log files. It allows automatic rotation, compression, removal, and mailing of log files. Each log file may be handled daily, weekly, monthly, or when it grows too large.

The logs are saved in `/var/log/watchdog/*.Log`, here is an example of a logrotate configuration
```
/var/log/watchdog/*.log {
    daily
    rotate 8
    size 500M
    compress
    delaycompress
    missingok
    notifempty
    create 440 root root
}
```