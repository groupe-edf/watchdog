package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"

	"github.com/groupe-edf/watchdog/cmd"
	"github.com/groupe-edf/watchdog/internal/util"
	logger "github.com/sirupsen/logrus"
)

var profile bool

func main() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := util.InitializeContext()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-interruptChan
		cancel()
	}()
	// Intercept application panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("something went worng, please contact your system administrator \n", string(debug.Stack()))
			os.Exit(0)
		}
	}()
	if profile {
		prof, err := os.Create("watchdog.prof")
		if err != nil {
			logger.Fatalf("could not create cpu profile: %v", err)
		}
		defer prof.Close()
		if err := pprof.StartCPUProfile(prof); err != nil {
			logger.Fatalf("could not start cpu profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}
	_ = cmd.Execute(ctx)
}
