package authentication

import (
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/server/authentication/provider"
	"github.com/groupe-edf/watchdog/internal/server/container"
)

const (
	// ServiceName service provider name
	ServiceName = "authentication"
)

// ServiceProvider authentication service provider
type ServiceProvider struct {
	Servers []*config.LDAP
}

// Register registring authentication service
func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(c container.Container) container.Service {
		authenticator := NewService(Options{
			BearerHeader: "Bearer",
			Secret:       []byte("secret"),
		})
		authenticator.AddProvider(provider.NewLDAPProvider(service.Servers))
		authenticator.AddProvider(provider.NewLocalProvider())
		return authenticator
	})
}
