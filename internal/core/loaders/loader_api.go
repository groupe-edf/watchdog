package loaders

import (
	"context"

	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/client"
)

type APILoader struct {
	client *client.Client
}

func (loader *APILoader) LoadPolicies(ctx context.Context) ([]models.Policy, error) {
	return loader.client.GetPolicies(ctx)
}

func (loader *APILoader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	return loader.client.GetRules(ctx)
}

func NewAPILoader(baseURL string, apiKey string) *APILoader {
	client := client.NewClient(baseURL, client.WithAPIKey(apiKey))
	return &APILoader{
		client: client,
	}
}
