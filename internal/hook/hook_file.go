package hook

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/version"
	"gopkg.in/yaml.v3"
)

var (
	// ConfigFilename Git hooks configuration filename to be parsed
	ConfigFilename = ".githooks.yml"
	// ErrFileNotFound File not found error
	ErrFileNotFound = errors.New("Configuration file .githooks.(yaml|yml) not found")
)

// LoadGitHooks load and return Configuration struct
func LoadGitHooks(filePath string) (*GitHooks, error) {
	content, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, ErrFileNotFound
	}
	return LoadGitHooksFromRaw(string(content))
}

// ExtractConfigFile extract .githooks.yml file from Git bare repository
func ExtractConfigFile(ctx context.Context, commit *object.Commit) (gitHooks *GitHooks, err error) {
	file, err := commit.File(ConfigFilename)
	if err != nil {
		return nil, ErrFileNotFound
	}
	content, err := file.Contents()
	if err != nil {
		return nil, err
	}
	return LoadGitHooksFromRaw(content)
}

// LoadGitHooksFromRaw load config file from raw data
func LoadGitHooksFromRaw(fileContent string) (*GitHooks, error) {
	var hooks = &GitHooks{}
	if err := yaml.Unmarshal([]byte(fileContent), hooks); err != nil {
		return nil, fmt.Errorf("Unable to decode into struct, %v", err)
	}
	if err := hooks.Validate(version.Version); err != nil {
		return nil, err
	}
	return hooks, nil
}
