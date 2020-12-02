package config

import (
	"os"
	"runtime"
	"strings"

	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/spf13/viper"
)

const (
	// LogsPath default log file location
	LogsPath = "/var/log/watchdog/watchdog.log"
)

// Options options data structure
type Options struct {
	AuthBasicToken     string               `mapstructure:"auth-basic-token"`
	Banner             bool                 `mapstructure:"banner"`
	CacheDirectory     string               `mapstructure:"cache-directory"`
	Contact            string               `mapstructure:"contact"`
	DefaultHandlers    map[string]hook.Rule `mapstructure:"default-handlers"`
	DocsLink           string               `mapstructure:"docs-link"`
	ErrorMessagePrefix string               `mapstructure:"error-message-prefix"`
	HookFile           string               `mapstructure:"hook-file"`
	HookInput          string               `mapstructure:"hook-input"`
	HookType           string               `mapstructure:"hook-type"`
	LogsFormat         string               `mapstructure:"logs-format"`
	LogsLevel          string               `mapstructure:"logs-level"`
	LogsPath           string               `mapstructure:"logs-path"`
	MaxFileSize        uint
	MaxRepositorySize  uint
	// MaxWorkers max workers running at the same time
	MaxWorkers       int    `mapstructure:"max-workers"`
	Output           string `mapstructure:"output"`
	OutputFormat     string `mapstructure:"output-format"`
	PluginsDirectory string `mapstructure:"plugins-directory"`
	Security         struct {
		MergeRules    bool `mapstructure:"merge-rules"`
		RevealSecrets int  `mapstructure:"reveal-secrets"`
		Rules         []struct {
			Description string   `mapstructure:"description"`
			File        string   `mapstructure:"file"`
			Regexp      string   `mapstructure:"regexp"`
			Severity    string   `mapstructure:"severity"`
			Tags        []string `mapstructure:"tags"`
		} `mapstructure:"rules"`
	} `mapstructure:"security"`
	Verbose bool   `mapstructure:"verbose"`
	URI     string `mapstructure:"uri"`
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
	if options.MaxWorkers == 0 {
		options.MaxWorkers = runtime.NumCPU()
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
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	config.SetEnvPrefix("WATCHDOG")
	config.AllowEmptyEnv(true)
	config.SetDefault("banner", true)
	err = config.Unmarshal(&options)
	if err != nil {
		return nil, err
	}
	err = config.UnmarshalKey("default_handlers", &options.DefaultHandlers)
	if err != nil {
		return nil, err
	}
	options.AuthBasicToken = config.GetString("auth_basic_token")
	options.CacheDirectory = config.GetString("cache_directory")
	options.ErrorMessagePrefix = config.GetString("error_message_prefix")
	options.Security.MergeRules = config.GetBool("security.merge_rules")
	options.Security.RevealSecrets = config.GetInt("security.reveal_secrets")
	err = options.Validate()
	if err != nil {
		return nil, err
	}
	return options, nil
}
