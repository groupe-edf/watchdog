package config

import (
	"os"
	"runtime"

	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/spf13/viper"
)

type DatabaseDriver string

const (
	PostgresDriver DatabaseDriver = "postgres"
	BoldDriver     DatabaseDriver = "bolt"
	ServiceName                   = "config"
)

type BoltOptions struct {
	Path string `mapstructure:"path"`
}
type PostgresOptions struct {
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
}

type Database struct {
	Bolt     BoltOptions     `mapstructure:"bolt"`
	Driver   DatabaseDriver  `mapstructure:"driver"`
	Postgres PostgresOptions `mapstructure:"postgres"`
}

type LDAP struct {
	Attributes    map[string]string `mapstructure:"attributes"`
	BindDN        string            `mapstructure:"bind_dn"`
	BindPassword  string            `mapstructure:"bind_password"`
	Label         string            `mapstructure:"label"`
	Host          string            `mapstructure:"host"`
	Port          int               `mapstructure:"port"`
	SearchBaseDNS string            `mapstructure:"search_base_dns"`
	SearchFilter  string            `mapstructure:"search_filter"`
	SSLSkipVerify bool              `mapstructure:"ssl_skip_verify"`
	StartSSL      bool              `mapstructure:"start_ssl"`
	UID           string            `mapstructure:"uid"`
	UseSSL        bool              `mapstructure:"use_ssl"`
}

type Logs struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
	Path   string `mapstructure:"path"`
}

type OAuthProvider struct {
	ApplicationID     string   `mapstructure:"application_id"`
	ApplicationSecret string   `mapstructure:"application_secret"`
	Args              []string `mapstructure:"args"`
	Name              string   `mapstructure:"name"`
}

// Options options data structure
type Options struct {
	*Logs              `mapstructure:"logs"`
	*Security          `mapstructure:"security"`
	*Server            `mapstructure:"server"`
	AuthBasicToken     string               `mapstructure:"auth-basic-token"`
	Banner             bool                 `mapstructure:"banner"`
	CacheDirectory     string               `mapstructure:"cache_directory"`
	Concurrent         int                  `mapstructure:"concurrent"` // Concurrent max workers running at the same time
	Contact            string               `mapstructure:"contact"`
	Handlers           map[string]hook.Rule `mapstructure:"handlers"`
	DocsLink           string               `mapstructure:"docs-link"`
	ErrorMessagePrefix string               `mapstructure:"error_message_prefix"`
	HookFile           string               `mapstructure:"hook-file"`
	HookInput          string               `mapstructure:"hook-input"`
	HookType           string               `mapstructure:"hook-type"`
	MaxFileSize        uint
	MaxRepositorySize  uint
	Output             string `mapstructure:"output"`
	OutputFormat       string `mapstructure:"output-format"`
	PluginsDirectory   string `mapstructure:"plugins-directory"`
	PoliciesFiles      string `mapstructure:"policies-file"`
	ServerURL          string `mapstructure:"server-url"`
	ServerToken        string `mapstructure:"server-token"`
	Verbose            bool   `mapstructure:"verbose"`
	URI                string `mapstructure:"uri"`
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
	Profile        bool             `mapstructure:"profile"`
	Database       *Database        `mapstructure:"database"`
	LDAP           []*LDAP          `mapstructure:"ldap"`
	ListenAddress  string           `mapstructure:"listen_address"`
	OAuthProviders []*OAuthProvider `mapstructure:"oauth_providers"`
	Security       struct {
		MasterKey string `mapstructure:"master_key"`
	} `mapstructure:"security"`
	Static struct {
		Routes []string `mapstructure:"routes"`
	} `mapstructure:"static"`
	Storage *Storage `mapstructure:"storage"`
}

type Storage struct {
	Directory string `mapstructure:"directory"`
	Driver    string `mapstructure:"driver"`
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
	err = config.UnmarshalKey("server", &options.Server)
	if err != nil {
		return nil, err
	}
	err = options.Validate()
	if err != nil {
		return nil, err
	}
	return options, nil
}
