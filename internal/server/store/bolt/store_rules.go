package bolt

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (store *BoltStore) FindRules(q *query.Query) (models.Paginator[models.Rule], error) {
	paginator := models.Paginator[models.Rule]{
		Items: make([]models.Rule, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) SaveRule(rule *models.Rule) (*models.Rule, error) {
	return nil, nil
}
func (store *BoltStore) ToggleRule(ruleID int64, enabled bool) error {
	return nil
}
