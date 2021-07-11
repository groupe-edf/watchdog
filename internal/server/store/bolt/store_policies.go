package bolt

import (
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (store *BoltStore) FindPolicies(q *query.Query) (models.Paginator[models.Policy], error) {
	paginator := models.Paginator[models.Policy]{
		Items: make([]models.Policy, 0),
		Query: q,
	}
	return paginator, nil
}
func (store *BoltStore) FindPolicyByID(id int64) (*models.Policy, error) {
	return nil, nil
}
func (store *BoltStore) TogglePolicy(id int64, enabled bool) error {
	return nil
}
