package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile  string
	rootCommand = &cobra.Command{
		Use:   "watchdog-server",
		Short: "",
		Long:  ``,
		PreRun: func(ccmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute execute audit command
func Execute(ctx context.Context) error {
	return rootCommand.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to watchdog configuration file")
	// Bind flags to configuration
	_ = viper.BindPFlags(rootCommand.PersistentFlags())
}

// Load configuration from file
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}
	viper.SetEnvPrefix("WATCHDOG_")
	viper.SetConfigType("yaml")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("/etc/watchdog/")
	viper.AddConfigPath("/etc/watchdog/config")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	for _, key := range viper.AllKeys() {
		viper.RegisterAlias(strings.ReplaceAll(key, "-", "_"), key)
	}
}
