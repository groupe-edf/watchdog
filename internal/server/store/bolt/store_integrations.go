package bolt

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (store *BoltStore) FindWebhooks(q *query.Query) (models.Paginator[models.Webhook], error) {
	paginator := models.Paginator[models.Webhook]{
		Items: make([]models.Webhook, 0),
		Query: q,
	}
	return paginator, nil
}
