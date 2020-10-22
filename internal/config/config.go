package config

import "os"

// Options options data structure
type Options struct {
	Banner           bool   `mapstructure:"banner"`
	Contact          string `mapstructure:"contact"`
	DocsLink         string `mapstructure:"docs-link"`
	ListenAddr       string `mapstructure:"listen-addr"`
	LogsFormat       string `mapstructure:"logs-format"`
	LogsLevel        string `mapstructure:"logs-level"`
	LogsPath         string `mapstructure:"logs-path"`
	HookFile         string `mapstructure:"hook-file"`
	HookInput        string `mapstructure:"hook-input"`
	PluginsDirectory string `mapstructure:"plugins-directory"`
	Verbose          bool   `mapstructure:"verbose"`
	URI              string `mapstructure:"uri"`
}

// Validate validate options
func (options *Options) Validate() error {
	if options.URI == "" {
		currentWorkingDirectory, _ := os.Getwd()
		options.URI = currentWorkingDirectory
	}
	return nil
}

// NewOptions return Options instance
func NewOptions() *Options {
	return &Options{
		Banner: true,
	}
}
