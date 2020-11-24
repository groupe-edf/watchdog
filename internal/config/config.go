package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	// LogsPath default log file location
	LogsPath = "/var/log/watchdog/watchdog.log"
)

// Options options data structure
type Options struct {
	Banner           bool   `mapstructure:"banner"`
	Contact          string `mapstructure:"contact"`
	DocsLink         string `mapstructure:"docs-link"`
	HookFile         string `mapstructure:"hook-file"`
	HookInput        string `mapstructure:"hook-input"`
	HookType         string `mapstructure:"hook-type"`
	LogsFormat       string `mapstructure:"logs-format"`
	LogsLevel        string `mapstructure:"logs-level"`
	LogsPath         string `mapstructure:"logs-path"`
	Output           string `mapstructure:"output"`
	OutputFormat     string `mapstructure:"output-format"`
	PluginsDirectory string `mapstructure:"plugins-directory"`
	Verbose          bool   `mapstructure:"verbose"`
	URI              string `mapstructure:"uri"`
}

// Validate validate options
func (options *Options) Validate() error {
	if options.URI == "" {
		directory, _ := os.Getwd()
		options.URI = directory
	}
	if options.LogsPath == "" {
		options.LogsPath = LogsPath
	}
	return nil
}

// NewOptions return Options instance
func NewOptions(config *viper.Viper) (options *Options, err error) {
	if err := config.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return options, err
		}
	}
	for _, key := range config.AllKeys() {
		config.RegisterAlias(strings.ReplaceAll(key, "-", "_"), key)
	}
	config.SetDefault("banner", true)
	err = config.Unmarshal(&options)
	if err != nil {
		return nil, err
	}
	err = options.Validate()
	if err != nil {
		return nil, err
	}
	return options, nil
}
