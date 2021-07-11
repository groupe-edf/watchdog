package loaders

import (
	"context"

	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

type StoreLoader struct {
	store models.Store
}

func (loader *StoreLoader) LoadPolicies(ctx context.Context) ([]models.Policy, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"policies"."enabled"`,
		Operator: query.Equal,
		Value:    true,
	})
	paginator, err := loader.store.FindPolicies(q)
	if err != nil {
		return nil, err
	}
	return paginator.Items, nil
}

func (loader *StoreLoader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"rules"."enabled"`,
		Operator: query.Equal,
		Value:    true,
	})
	paginator, err := loader.store.FindRules(q)
	if err != nil {
		return nil, err
	}
	return paginator.Items, nil
}

func NewStoreLoader(store models.Store) *StoreLoader {
	return &StoreLoader{
		store: store,
	}
}
