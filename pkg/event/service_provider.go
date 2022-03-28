package event

import "github.com/groupe-edf/watchdog/pkg/container"

const ServiceName = "event_bus"

type ServiceProvider struct {
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(c container.Container) container.Service {
		return NewEventBus()
	})
}
