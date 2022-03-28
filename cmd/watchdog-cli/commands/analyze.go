package commands

import (
	"context"
	"os"
	"os/signal"

	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/core/loaders"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/output"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/logging"
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
			di.Provide(&git.ServiceProvider{
				Options: options,
			})
			di.Provide(&logging.ServiceProvider{
				Options: options.Logs,
			})
			di.Set(config.ServiceName, func(c container.Container) container.Service {
				return options
			})
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
			driver := di.Get(git.ServiceName).(git.Driver)
			repository, err := driver.Clone(ctx, git.CloneOptions{
				URL: options.URI,
			})
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
			}).Debugf("repository `%v` successfully fetched in `%s`", options.URI, repository.Path())
			loader, _ := loaders.GetLoader(options)
			analyzer, err := core.NewAnalyzer(ctx, loader, logger, options, repository, models.Whitelist{})
			if err != nil {
				logger.Fatal(err)
			}
			analyzer.SetDriver(driver)
			commitIter, err := driver.Commits(ctx, git.LogOptions{})
			if err != nil {
				logger.Fatal(err)
			}
			analyzeChan := make(chan models.AnalysisResult)
			go analyzer.Analyze(repository, commitIter, analyzeChan)
			writer := output.NewConsole(analyzeChan)
			writer.WriteTo()
		},
	}
)

func init() {
	rootCommand.AddCommand(analyzeCommand)
}
