package loaders

import (
	"context"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/git"
)

type GitLoader struct {
	backend    git.Driver
	repository *git.Repository
}

func (loader *GitLoader) LoadPolicies(ctx context.Context) ([]models.Policy, error) {
	fileContent, err := loader.backend.File(ConfigFilename)
	if err != nil {
		return nil, err
	}
	policies := make([]models.Policy, 0)
	return load(fileContent, policies)
}

func (loader *GitLoader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	fileContent, err := loader.backend.File(ConfigFilename)
	if err != nil {
		return nil, err
	}
	rules := make([]models.Rule, 0)
	return load(fileContent, rules)
}

func NewGitLoader(backend git.Driver) *GitLoader {
	return &GitLoader{
		backend: backend,
	}
}
