package notification

import (
	"github.com/groupe-edf/watchdog/pkg/container"
)

const (
	ServiceName = "notification"
)

type ServiceProvider struct {
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		return NewManager()
	})
}
