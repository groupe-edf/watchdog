package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/output"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/groupe-edf/watchdog/pkg/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	options        *config.Options
	analyzeCommand = &cobra.Command{
		Use:   "analyze",
		Short: "Run analysis",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var ctx = cmd.Context()
			var err error
			var exitMessage = "Your push was successfully accepted"
			var exitStatus = 0
			var hooks *hook.GitHooks
			var info *hook.Info
			options, err = config.NewOptions(viper.GetViper())
			if err != nil {
				color.Red.Printf("unable to decode into config struct, %v", err)
				os.Exit(0)
			}
			if options.Banner {
				_ = util.PrintBanner(ctx, options)
				fmt.Println()
			}
			logger := logging.New(logging.Options{
				LogsFormat: options.LogsFormat,
				LogsLevel:  options.LogsLevel,
				LogsPath:   options.LogsPath,
			})
			ctx, cancel := context.WithCancel(ctx)
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
			analyzer, err := core.NewAnalyzer(nil, options)
			if err != nil {
				logger.Fatal(err)
			}
			analyzer.SetLogger(logger)
			// Loading git repository
			repository, err := util.LoadRepository(ctx, options)
			logger.WithFields(logging.Fields{
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Debugf("repository `%v` successfully fetched", options.URI)
			if err != nil {
				color.Red.Printf("error fetching repository `%v`", err)
				logger.WithFields(logging.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Fatalf("error fetching repository `%v`", err)
			}
			analyzer.SetRepository(repository)
			// Loading .githooks.yml file
			if options.HookFile != "" {
				// External .githooks.yml
				logger.WithFields(logging.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Debugf("loading external hooks file %v", options.HookFile)
				hooks, err = hook.LoadGitHooks(options.HookFile)
				if err != nil {
					color.Red.Printf("error loading git hooks %v", err)
					os.Exit(0)
				}
			} else {
				// Versionned .githooks.yml
				var commit *object.Commit
				info, err = hook.ParseInfo(repository, options.HookInput)
				if err != nil && err != hook.ErrNoHookData {
					color.Red.Printf("error parsing hook info %v", err)
					logger.WithFields(logging.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
				}
				if info != nil {
					analyzer.SetInfo(info)
					commit = info.NewRev
				} else {
					reference, err := repository.Head()
					if err != nil {
						logger.WithFields(logging.Fields{
							"correlation_id": util.GetRequestID(ctx),
							"user_id":        util.GetUserID(ctx),
						}).Fatal(err)
					}
					commit, err = repository.CommitObject(reference.Hash())
					if err != nil {
						color.Red.Println(err.Error())
						logger.WithFields(logging.Fields{
							"commit":         commit.Hash.String(),
							"correlation_id": util.GetRequestID(ctx),
							"user_id":        util.GetUserID(ctx),
						}).Fatal(err)
					}
				}
				hooks, err = hook.ExtractConfigFile(ctx, commit)
				if err != nil && !errors.Is(err, hook.ErrFileNotFound) {
					color.Red.Println(err.Error())
					logger.WithFields(logging.Fields{
						"commit":         commit.Hash.String(),
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatalf("error when extracting config file %v", err)
				}
			}
			// No .githooks.yml file was referenced, create default one if we have global default handlers in configuration
			if hooks == nil && len(options.Handlers) > 0 {
				hooks = &hook.GitHooks{
					Hooks: []hook.Hook{
						{
							Name: "default",
						},
					},
					Version: version.Version,
				}
				logger.WithFields(logging.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Debug("no .githooks.yml file was referenced")
			}
			if hooks != nil {
				analyzer.SetHooks(hooks)
				// Register handlers
				analyzer.RegisterHandler(ctx, &handlers.BranchHandler{})
				analyzer.RegisterHandler(ctx, &handlers.CommitHandler{})
				analyzer.RegisterHandler(ctx, &handlers.FileHandler{})
				analyzer.RegisterHandler(ctx, &handlers.JiraHandler{})
				analyzer.RegisterHandler(ctx, &handlers.SecurityHandler{})
				analyzer.RegisterHandler(ctx, &handlers.TagHandler{})
				// Fetching commits
				commits, err := util.FetchCommits(repository, info, options.HookType)
				fmt.Println()
				if err != nil {
					color.Red.Println(err.Error())
					logger.WithFields(logging.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
				}
				// Run analysis
				err = analyzer.Analyze(ctx, commits)
				if err != nil {
					logger.WithFields(logging.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
				}
				// Send report
				err = output.Report(options, analyzer.Issues)
				if err != nil {
					logger.WithFields(logging.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
					os.Exit(0)
				}
				if analyzer.HasErrors() {
					exitMessage = "Your push was rejected because previous errors"
					exitStatus = 1
				}
				if info != nil {
					util.PrintMessage(exitMessage)
				}
				fmt.Println()
				util.ElapsedTime(ctx, "Operation")
			}
			os.Exit(exitStatus)
		},
	}
)

func init() {
	rootCommand.AddCommand(analyzeCommand)
}
