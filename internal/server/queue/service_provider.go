package queue

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

const (
	ServiceName = "queue"
)

type ServiceProvider struct {
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		logger := di.Get(logging.ServiceName).(logging.Interface)
		store := di.Get(store.ServiceName).(models.Store)
		return NewClient(store, logger)
	})
}
