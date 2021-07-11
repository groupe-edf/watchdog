package config

import (
	"os"
	"runtime"

	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/spf13/viper"
)

type Database struct {
	Host     string
	Name     string
	Password string
	Port     int
	Username string
}

// Options options data structure
type Options struct {
	AuthBasicToken     string `mapstructure:"auth-basic-token"`
	Banner             bool   `mapstructure:"banner"`
	CacheDirectory     string `mapstructure:"cache-directory"`
	Contact            string `mapstructure:"contact"`
	Database           Database
	Handlers           map[string]hook.Rule `mapstructure:"handlers"`
	DocsLink           string               `mapstructure:"docs-link"`
	ErrorMessagePrefix string               `mapstructure:"error_message_prefix"`
	HookFile           string               `mapstructure:"hook-file"`
	HookInput          string               `mapstructure:"hook-input"`
	HookType           string               `mapstructure:"hook-type"`
	LogsFormat         string               `mapstructure:"logs-format"`
	LogsLevel          string               `mapstructure:"logs-level"`
	LogsPath           string               `mapstructure:"logs-path"`
	MaxFileSize        uint
	MaxRepositorySize  uint
	// Concurrent max workers running at the same time
	Concurrent       int    `mapstructure:"concurrent"`
	Output           string `mapstructure:"output"`
	OutputFormat     string `mapstructure:"output-format"`
	PluginsDirectory string `mapstructure:"plugins-directory"`
	*Security        `mapstructure:"security"`
	*Server					 `mapstructure:"server"`
	Verbose          bool   `mapstructure:"verbose"`
	URI              string `mapstructure:"uri"`
}

// Security settings
type Security struct {
	MergeRules    bool `mapstructure:"merge_rules"`
	RevealSecrets int  `mapstructure:"reveal_secrets"`
	Rules         []struct {
		Description string   `mapstructure:"description"`
		File        string   `mapstructure:"file"`
		Regexp      string   `mapstructure:"regexp"`
		Severity    string   `mapstructure:"severity"`
		Tags        []string `mapstructure:"tags"`
	} `mapstructure:"rules"`
}

// Server settings
type Server struct {
	LDAP struct {
		BindDN string   `mapstructure:"bind_db"`
		BindPassword string   `mapstructure:"bind_password"`
		Host string   `mapstructure:"host"`
		Port int `mapstructure:"port"`
		SearchBaseDNS string   `mapstructure:"search_base_dns"`
		SearchFilter string   `mapstructure:"search_filter"`
		SSLSkipVerify bool `mapstructure:"ssl_skip_verify"`
		StartSSL bool `mapstructure:"start_ssl"`
		UseSSL bool `mapstructure:"use_ssl"`
	} `mapstructure:"ldap"`
	ListenAddress string `mapstructure:"listen_address"`
}

// Validate validate options
func (options *Options) Validate() error {
	if options.URI == "" {
		directory, _ := os.Getwd()
		options.URI = directory
	}
	if options.Concurrent == 0 {
		options.Concurrent = runtime.NumCPU()
	}
	if options.Security != nil && len(options.Security.Rules) > 0 {
		for index, rule := range options.Security.Rules {
			if rule.Description == "" || rule.Regexp == "" {
				options.Security.Rules = append(options.Security.Rules[:index], options.Security.Rules[index+1:]...)
			}
			if rule.Severity == "" {
				rule.Severity = "INFO"
			}
		}
	}
	return nil
}

// NewOptions return Options instance
func NewOptions(config *viper.Viper) (options *Options, err error) {
	err = config.Unmarshal(&options)
	if err != nil {
		return nil, err
	}
	err = config.UnmarshalKey("handlers", &options.Handlers)
	if err != nil {
		return nil, err
	}
	err = config.UnmarshalKey("security", &options.Security)
	if err != nil {
		return nil, err
	}
	err = options.Validate()
	if err != nil {
		return nil, err
	}
	return options, nil
}
