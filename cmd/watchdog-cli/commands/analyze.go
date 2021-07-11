package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-git/go-git/v5"
	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/backend"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/core/loaders"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/output"
	"github.com/groupe-edf/watchdog/internal/server/container"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	analyzeCommand = &cobra.Command{
		Use:   "analyze",
		Short: "Run analysis",
		Long:  ``,
		Run: func(cmd *cobra.Command, _ []string) {
			options, err := config.NewOptions(viper.GetViper())
			if err != nil {
				color.Red.Printf("unable to decode into config struct, %v", err)
				os.Exit(0)
			}
			di := container.GetContainer()
			di.Provide(&logging.ServiceProvider{
				Options: options.Logs,
			})
			di.Set(config.ServiceName, func(c container.Container) container.Service {
				return options
			})
			if options.Banner {
				_ = util.PrintBanner(options)
				fmt.Println()
			}
			logger := di.Get(logging.ServiceName).(logging.Interface)
			ctx, cancel := context.WithCancel(cmd.Context())
			interruption := make(chan os.Signal, 1)
			signal.Notify(interruption, os.Interrupt)
			defer func() {
				signal.Stop(interruption)
				cancel()
			}()
			go func() {
				select {
				case <-interruption:
					cancel()
				case <-ctx.Done():
				}
			}()
			client := backend.New(options)
			err = client.Clone(ctx, &git.CloneOptions{
				Progress: os.Stderr,
				URL:      options.URI,
			})
			if err != nil {
				logger.Fatal(err)
			}
			repository, err := util.LoadRepository(ctx, options)
			if err != nil {
				color.Red.Printf("error fetching repository `%v`", err)
				logger.WithFields(logging.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Fatalf("error fetching repository `%v`", err)
			}
			logger.WithFields(logging.Fields{
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Debugf("repository `%v` successfully fetched", options.URI)
			loader := loaders.NewAPILoader(options.ServerURL, options.ServerToken)
			analyzer, err := core.NewAnalyzer(ctx, loader, logger, options, models.Whitelist{})
			if err != nil {
				logger.Fatal(err)
			}
			analyzer.SetRepository(repository)
			commitIter, err := client.Commits(ctx, &git.LogOptions{})
			if err != nil {
				logger.Fatal(err)
			}
			defer commitIter.Close()
			analyzeChan := make(chan models.AnalysisResult)
			go analyzer.Analyze(commitIter, analyzeChan)
			writer := output.NewConsole(analyzeChan)
			writer.WriteTo()
		},
	}
)

func init() {
	rootCommand.AddCommand(analyzeCommand)
}
