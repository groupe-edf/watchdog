package job

import (
	"context"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

type ProcessAnalytics struct {
	Context context.Context
	Logger  logging.Interface
	Options *config.Options
	Store   models.Store
}

func (processor *ProcessAnalytics) Handle(job *models.Job) error {
	err := processor.Store.RefreshAnalytics()
	if err != nil {
		return err
	}
	return nil
}
