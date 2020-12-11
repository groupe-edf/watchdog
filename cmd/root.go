package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFile  string
	rootCommand = &cobra.Command{
		Use:   "watchdog",
		Short: "",
		Long:  ``,
	}
)

// Execute execute audit command
func Execute(ctx context.Context) error {
	return rootCommand.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCommand.PersistentFlags().Bool("profile", false, "collect the profile to hercules.pprof.")
	rootCommand.PersistentFlags().Int("concurrent", 0, "concurrent worker used to run analysus")
	rootCommand.PersistentFlags().Int("security.reveal-secrets", 0, "full or partial reveal of secrets in report and logs")
	rootCommand.PersistentFlags().String("auth-basic-token", "", "authentication token used to fetch remote repositories")
	rootCommand.PersistentFlags().String("hook-input", "", "standard input <old-value> SP <new-value> SP <ref-name> LF")
	rootCommand.PersistentFlags().String("hook-type", "", "git server-side hook pre-receive, update or post-receive")
	rootCommand.PersistentFlags().String("docs-link", "", "link to documentation")
	rootCommand.PersistentFlags().String("logs-format", "json", "logging level")
	rootCommand.PersistentFlags().String("logs-level", "info", "logging level")
	rootCommand.PersistentFlags().String("logs-path", "/var/log/watchdog/watchdog.log", "path to logs")
	rootCommand.PersistentFlags().String("output", "", "path to output file")
	rootCommand.PersistentFlags().String("output-format", "text", "report format")
	rootCommand.PersistentFlags().String("plugins-directory", "plugins", "path to plugins directory")
	rootCommand.PersistentFlags().String("uri", "", "path to working directory")
	rootCommand.PersistentFlags().StringP("hook-file", "f", "", "path to external .githooks.yml file")
	rootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to watchdog configuration file")
	rootCommand.PersistentFlags().Bool("verbose", true, "make the operation more talkative")
	// Bind flags to configuration
	_ = viper.BindPFlags(rootCommand.PersistentFlags())
}

// Load configuration from file
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/watchdog/")
	viper.AddConfigPath("/etc/watchdog/config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("WATCHDOG")
	viper.AllowEmptyEnv(true)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	for _, key := range viper.AllKeys() {
		viper.RegisterAlias(strings.ReplaceAll(key, "-", "_"), key)
	}
	rootCommand.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && viper.IsSet(f.Name) {
			value := viper.Get(f.Name)
			_ = rootCommand.Flags().Set(f.Name, fmt.Sprintf("%v", value))
		}
	})
}
