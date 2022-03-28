package bolt

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (store *BoltStore) AddPolicyCondition(condition *models.Condition) (*models.Condition, error) {
	return nil, nil
}
func (store *BoltStore) DeletePolicy(policyID int64) error {
	return nil
}
func (store *BoltStore) DeletePolicyCondition(policyID int64, conditionID int64) error {
	return nil
}
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
func (store *BoltStore) SavePolicy(policy *models.Policy) (*models.Policy, error) {
	return nil, nil
}
func (store *BoltStore) TogglePolicy(id int64, enabled bool) error {
	return nil
}
func (store *BoltStore) UpdatePolicy(policy *models.Policy) (*models.Policy, error) {
	return nil, nil
}
