---
title: "Install"
linkTitle: "Install"
weight: 2
description: >
  How to download and install Watchdog.
---

## From the Binary Releases
Every [release](https://github.com/groupe-edf/watchdog/releases) of Watchdog provides binary releases for a variety of OSes. These binary versions can be manually downloaded and installed.

* Download your desired version
* Unpack it `tar -zxvf watchdog-v1.2.0-linux-amd64.tar.gz`
* Find the `watchdog` binary in the unpacked directory, and move it to its desired destination `mv watchdog-v1.2.0-linux-amd64 /usr/local/bin/watchdog`
* Make it executable `chmod +x /usr/local/bin/watchdog`


```bash
watchdog version
```