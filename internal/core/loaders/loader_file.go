package loaders

import (
	"context"
	"io/fs"
	"os"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

type FileLoader struct {
	fileSystem fs.FS
}

func (loader *FileLoader) LoadPolicies(ctx context.Context) ([]models.Policy, error) {
	fileContent, err := loader.loadFile()
	if err != nil {
		return nil, err
	}
	policies := make([]models.Policy, 0)
	return load(string(fileContent), policies)
}

func (loader *FileLoader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	fileContent, err := loader.loadFile()
	if err != nil {
		return nil, err
	}
	rules := make([]models.Rule, 0)
	return load(string(fileContent), rules)
}

func (loader *FileLoader) loadFile() ([]byte, error) {
	matches, err := fs.Glob(loader.fileSystem, ConfigFilename)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, ErrFileNotFound
	}
	return fs.ReadFile(loader.fileSystem, matches[0])
}

func NewFileLoader(path string) *FileLoader {
	fileSystem := os.DirFS(path)
	return &FileLoader{
		fileSystem: fileSystem,
	}
}
