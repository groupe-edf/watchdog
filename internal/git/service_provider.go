package git

import (
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/pkg/container"
)

const ServiceName = "git"

type ServiceProvider struct {
	Options *config.Options
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		driver := NewGit(service.Options)
		return driver
	})
}
