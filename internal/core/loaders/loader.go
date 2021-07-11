package loaders

import (
	"context"
	"errors"
	"fmt"

	"github.com/groupe-edf/watchdog/internal/models"
	"gopkg.in/yaml.v3"
)

var (
	// ConfigFilename Git hooks configuration filename to be parsed
	ConfigFilename = ".watchdog.*"
	// ErrFileNotFound File not found error
	ErrFileNotFound = errors.New("configuration file .githooks.(yaml|yml) not found")
)

type Loader interface {
	LoadPolicies(ctx context.Context) ([]models.Policy, error)
	LoadRules(ctx context.Context) ([]models.Rule, error)
}

func load[T any](fileContent string, data T) (T, error) {
	if err := yaml.Unmarshal([]byte(fileContent), &data); err != nil {
		return data, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return data, nil
}
