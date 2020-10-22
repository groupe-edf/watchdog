---
title: "Release"
linkTitle: "Release"
weight: 1
description: >
  Release binaries and docker image.
---

Docker
--------------------
Create the Watchdog Docker image with the following command.
``` bash
docker build \
  --tag groupe-edf/watchdog \
  --file build/Dockerfile \
  --build-arg=http_proxy=$HTTPS_PROXY \
  --build-arg=https_proxy=$HTTPS_PROXY \
  --build-arg=CGO_ENABLED=0 \
  .
```