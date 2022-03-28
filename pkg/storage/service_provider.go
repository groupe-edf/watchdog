package storage

import "github.com/groupe-edf/watchdog/pkg/container"

const (
	// ServiceName service provider name
	ServiceName = "storage"
)

// ServiceProvider storage service provider
type ServiceProvider struct {
}

// Register registring storage service
func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(c container.Container) container.Service {
		return &Local{}
	})
}
