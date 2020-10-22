---
title: "Bootstrap"
linkTitle: "Bootstrap"
weight: 2
description: >
  How to setup Watchdog.
---

Developing Watchdog
--------------------

If you wish to work on Watchdog itself, you'll first need Go installed on your machine.
Go version 1.15.2+ is required.

For local dev first make sure Go is properly installed, including setting up a
[GOPATH](https://golang.org/doc/code.html#GOPATH). Ensure that `$GOPATH/bin` is in
your path as some distributions bundle old version of build tools. Next, clone this
repository. Watchdog uses [Go Modules](https://github.com/golang/go/wiki/Modules),
so it is recommended that you clone the repository ***outside*** of the GOPATH.
You can then download any required build tools by bootstrapping your environment:

```sh
$ make bootstrap
```

## Prerequisite
* Clone the source code
```bash
$ git clone https://github.com/groupe-edf/watchdog
```
Now you can build and run Watchdog by one of the following ways

## Build and run Watchdog locally
1. Build Watchdog binary
```bash
# Fetch the dependencies
$ go mod download
# Build the binary
$ go build -o watchdog
```
2. Run Watchdog binary
```bash
./watchdog version
```

## Making A Change
* Before making any significant changes, please [open an issue](https://github.com/groupe-edf/watchdog/issues). Discussing your proposed changes ahead of time will make the contribution process smooth for everyone.
* Once we’ve discussed your changes and you’ve got your code ready, make sure that build steps pass. Open your pull request against `develop` branch.
* To avoid build failures in CI, run
```bash
$ make lint
$ make test-unit
```
This will check if the code is properly formatted, linted.
* Run security and e2e tests
```bash
$ make test-security
$ make test-integration
```
* Make sure your pull request has [good commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
* Try to squash unimportant commits and rebase your changes on to develop branch, this will make sure we have clean log of changes.