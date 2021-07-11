package event

import "github.com/groupe-edf/watchdog/internal/server/container"

const ServiceName = "event_manager"

type ServiceProvider struct {
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(c container.Container) container.Service {
		return NewManager()
	})
}
