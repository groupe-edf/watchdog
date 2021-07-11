package loaders

import (
	"context"
	"io"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/groupe-edf/watchdog/internal/models"
)

type GitLoader struct {
	repository *git.Repository
}

func (loader *GitLoader) LoadPolicies(ctx context.Context) ([]models.Policy, error) {
	fileContent, err := loader.loadFile()
	if err != nil {
		return nil, err
	}
	policies := make([]models.Policy, 0)
	return load(fileContent, policies)
}

func (loader *GitLoader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	fileContent, err := loader.loadFile()
	if err != nil {
		return nil, err
	}
	rules := make([]models.Rule, 0)
	return load(fileContent, rules)
}

func (loader *GitLoader) loadFile() (string, error) {
	worktree, err := loader.repository.Worktree()
	if err != nil {
		return "", err
	}
	file, err := worktree.Filesystem.Open(".watchdog.yml")
	if err != nil {
		return "", err
	}
	defer file.Close()
	builder := &strings.Builder{}
	_, err = io.Copy(builder, file)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func NewGitLoader(repository *git.Repository) *GitLoader {
	return &GitLoader{
		repository: repository,
	}
}
