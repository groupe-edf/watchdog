---
title: "Install"
linkTitle: "Install"
weight: 1
description: >
  How to setup your development environment
---

## Downloading Go Lang Installation Files
Visit Go Lang’s official download page to download the installer according to your operating system. See [Downloads](https://golang.org/dl/)

> Go Lang is currently supported only on **x64** based processors.

### Setting up on Windows
* Double click the MSI installer, you just downloaded to start the Go Lang installation on your Windows system.
* Follow the prompts to install the Go tools. By default, the installer puts the Go distribution in `C:\Go`.
* The installer should put the `C:\Go\bin` directory in your `PATH` environment variable. You may need to restart any open command prompts for the change to take effect.

### Setting up on Linux
* Extract the `.tar.gz` archive you just downloaded to `/usr/local`, creating a Go tree in `/usr/local/go`.
```bash
$ sudo tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz
```
> Change the archive name as per requirement. go1.15.2.linux-amd64.tar.gz is the latest version available at the time of this writing.
* Add `/usr/local/go/bin` to the PATH environment variable. You can do this by adding this line to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:
```bash
$ export PATH=$PATH:/usr/local/go/bin
```

## Verifying your installation
```bash
$ go version
```
* Create a file named `main.go` using the text editor of your choice, and copy-paste the following content into it.
```golang
package main

import "fmt"

func main() {
    fmt.Printf("Hello, World..!!\n")
}
```
* Save the file and go back to the terminal and run the following command to compile the code into a binary.
```bash
$ go build
```
* The previous step creates a binary named `main` in the same directory. Simply run the executable to run the code.
* If the installation had been correct, the exe will print `Hello, World..!!`