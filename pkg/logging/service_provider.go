package logging

import (
	"os"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/pkg/container"
)

const ServiceName = "logger"

type ServiceProvider struct {
	Options *config.Logs
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		logger := New(Options{
			LogsFormat:       service.Options.Format,
			LogsLevel:        service.Options.Level,
			LogsOutput:       os.Stdout,
			LogsPath:         service.Options.Path,
			LogsReportCaller: true,
		})
		return logger
	})
}
