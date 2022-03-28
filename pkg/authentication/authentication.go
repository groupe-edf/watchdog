package authentication

import (
	"github.com/groupe-edf/watchdog/pkg/authentication/provider"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

// Options authentication service settings
type Options struct {
	BearerHeader  string
	Secret        []byte
	JWTCookieName string // default "JWT"
	JWTHeaderKey  string // default "Authorization"
	JWTQueryParam string // default "authorization"
	Issuer        string // default "https://carecrute.com"
}

// Authentication authentication service provider
type Authentication struct {
	Providers []provider.Provider
	Options   Options
}

func (authentication *Authentication) AddProvider(provider provider.Provider) {
	authentication.Providers = append(authentication.Providers, provider)
}

func (authentication *Authentication) Authenticate(email string, password string) (identity *provider.Identity, err error) {
	logger := container.Get(logging.ServiceName).(logging.Interface)
	for _, provider := range authentication.Providers {
		identity, err = provider.Authenticate(email, password)
		if err == nil && identity != nil {
			return identity, err
		}
		if err != nil {
			logger.Error(err)
		}
	}
	return identity, err
}

func NewService(options Options) *Authentication {
	if options.JWTQueryParam == "" {
		options.JWTQueryParam = "authorization"
	}
	if options.JWTHeaderKey == "" {
		options.JWTHeaderKey = "Authorization"
	}
	return &Authentication{
		Options: options,
	}
}
