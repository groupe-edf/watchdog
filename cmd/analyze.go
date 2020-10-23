package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/output"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile     string
	hookType       string
	options        = config.NewOptions()
	analyzeCommand = &cobra.Command{
		Use:   "analyze",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var ctx = cmd.Context()
			var err error
			var exitMessage = "Your push was successfully accepted"
			var exitStatus = 0
			var hooks *hook.GitHooks
			var info *hook.Info
			err = viper.Unmarshal(options)
			if err != nil {
				fmt.Printf("Unable to decode into config struct, %v", err)
				os.Exit(0)
			}
			if err := options.Validate(); err != nil {
				fmt.Printf("Unable to decode into config struct, %v", err)
				os.Exit(0)
			}
			if options.Banner {
				_ = util.PrintBanner(ctx, options)
				fmt.Println()
			}
			logger := util.GetLogger(options)
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
			logger.WithFields(logrus.Fields{
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Debugf("Loading repository `%v`", options.URI)
			repository, err := util.LoadRepository(options.URI)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Fatal(err)
			}
			analyzer.SetRepository(repository)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Warning(err)
			}
			if options.HookFile != "" {
				logger.WithFields(logrus.Fields{
					"correlation_id": util.GetRequestID(ctx),
					"user_id":        util.GetUserID(ctx),
				}).Debugf("loading external hooks file %v", options.HookFile)
				hooks, err = hook.LoadGitHooks(options.HookFile)
				if err != nil {
					fmt.Printf("Error loading git hooks %v", err)
					os.Exit(0)
				}
			} else {
				var commit *object.Commit
				info, err = hook.ParseInfo(repository, options.HookInput)
				if err != nil {
					fmt.Printf("Error parsing hook info %v", err)
					os.Exit(0)
				}
				if info != nil {
					analyzer.SetInfo(info)
					commit = info.NewRev
				} else {
					reference, err := repository.Head()
					if err != nil {
						logger.WithFields(logrus.Fields{
							"correlation_id": util.GetRequestID(ctx),
							"user_id":        util.GetUserID(ctx),
						}).Fatal(err)
					}
					commit, err = repository.CommitObject(reference.Hash())
					if err != nil {
						logger.WithFields(logrus.Fields{
							"commit":         commit.Hash.String(),
							"correlation_id": util.GetRequestID(ctx),
							"user_id":        util.GetUserID(ctx),
						}).Fatal(err)
					}
				}
				hooks, err = hook.ExtractConfigFile(ctx, commit)
				if err != nil && !errors.Is(err, hook.ErrFileNotFound) {
					logger.WithFields(logrus.Fields{
						"commit":         commit.Hash.String(),
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Errorf("Error when extracting config file %v", err)
				}
			}
			if hooks != nil {
				analyzer.SetHooks(hooks)
				// Register handlers
				analyzer.RegisterHandler(ctx, &handlers.BranchHandler{})
				analyzer.RegisterHandler(ctx, &handlers.CommitHandler{})
				analyzer.RegisterHandler(ctx, &handlers.FileHandler{})
				analyzer.RegisterHandler(ctx, &handlers.JiraHandler{})
				analyzer.RegisterHandler(ctx, &handlers.SecurityHandler{})
				// Fetching commits
				commits, err := util.FetchCommits(repository, info, hookType)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
				}
				if len(commits) == 0 {
					logger.WithFields(logrus.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(errors.New("No commits found"))
				}
				err = analyzer.Analyze(ctx, commits)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
				}
				err = output.Report(viper.GetString("output"), viper.GetString("output-format"), analyzer.Issues)
				if err != nil {
					logger.WithFields(logrus.Fields{
						"correlation_id": util.GetRequestID(ctx),
						"user_id":        util.GetUserID(ctx),
					}).Fatal(err)
					os.Exit(0)
				}
				if analyzer.HasErrors() {
					exitMessage = "Your push was rejected because previous errors"
					exitStatus = 1
				}
				util.PrintMessage(exitMessage)
				util.ElapsedTime(ctx, "Operation")
			}
			os.Exit(exitStatus)
		},
	}
)

// Execute execute audit command
func Execute(ctx context.Context) error {
	return analyzeCommand.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)
	analyzeCommand.Flags().Bool("profile", false, "collect the profile to hercules.pprof.")
	analyzeCommand.PersistentFlags().String("hook-input", "", "standard input <old-value> SP <new-value> SP <ref-name> LF")
	analyzeCommand.Flags().StringVar(&hookType, "hook-type", "pre-receive", "git server-side hook pre-receive, update or post-receive")
	analyzeCommand.PersistentFlags().String("docs-link", "", "link to documentation")
	analyzeCommand.PersistentFlags().String("logs-format", "json", "logging level")
	analyzeCommand.PersistentFlags().String("logs-level", "info", "logging level")
	analyzeCommand.PersistentFlags().String("logs-path", "/var/log/watchdog/watchdog.log", "path to logs")
	analyzeCommand.PersistentFlags().String("output", "", "path to output file")
	analyzeCommand.PersistentFlags().String("output-format", "text", "report format")
	analyzeCommand.PersistentFlags().String("plugins-directory", "plugins", "path to plugins directory")
	analyzeCommand.PersistentFlags().String("uri", "", "path to working directory")
	analyzeCommand.PersistentFlags().StringP("hook-file", "f", "", "path to external .githooks.yml file")
	analyzeCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to watchdog configuration file")
	analyzeCommand.PersistentFlags().Bool("verbose", true, "make the operation more talkative")
	_ = analyzeCommand.MarkFlagRequired("hook-input")
	_ = analyzeCommand.MarkFlagRequired("hook-type")
	// Bind flags to configuration
	_ = viper.BindPFlags(analyzeCommand.PersistentFlags())
}

// Load configuration from file
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.AutomaticEnv()
	}
}
