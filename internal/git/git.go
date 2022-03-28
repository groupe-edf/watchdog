package git

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/models"
)

var _ Driver = (*Git)(nil)

type Git struct {
	ctx        context.Context
	options    *config.Options
	repository *Repository
}

func (driver *Git) Clone(ctx context.Context, cloneOptions CloneOptions) (*Repository, error) {
	if cloneOptions.URL == "" {
		return nil, errors.New("url is mandatory")
	}
	var repository *Repository
	if strings.Contains(cloneOptions.URL, "://") || regexp.MustCompile(string(`^[A-Za-z]\w*@[A-Za-z0-9][\w.]*:`)).MatchString(cloneOptions.URL) {
		hasher := sha1.New()
		hasher.Write([]byte(cloneOptions.URL))
		cachePath := filepath.Join(driver.options.CacheDirectory, hex.EncodeToString(hasher.Sum(nil)))
		err := os.RemoveAll(cachePath)
		if err != nil {
			return nil, err
		}
		if err := Clone(driver.ctx, cloneOptions.URL, cachePath, cloneOptions); err != nil {
			_ = os.RemoveAll(cachePath)
			return nil, err
		}
		repository, _ = NewRepository(cachePath)
	} else {
		repository, _ = NewRepository(cloneOptions.URL)
	}
	driver.repository = repository
	return driver.repository, nil
}

func (driver *Git) CommitsCount(revision string) (int64, error) {
	return 0, nil
}

func (driver *Git) Diff(startCommit, endCommit string) error {
	return nil
}

func (driver *Git) File(filePath string) (string, error) {
	return "", nil
}

func (driver *Git) GetRepository() *Repository {
	return driver.repository
}

func (driver *Git) Head() (string, error) {
	output, _, err := NewCommand(context.Background(), "rev-parse", "HEAD").RunStdString(&RunOptions{
		Dir: driver.repository.Path(),
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(output, "\n"), nil
}

func (driver *Git) Commits(ctx context.Context, options LogOptions) (models.Iterator[models.Commit], error) {
	output, _, err := NewCommand(context.Background(), "log", "--pretty=format:%H").RunStdBytes(&RunOptions{
		Dir: driver.repository.Path(),
	})
	if err != nil {
		return nil, err
	}
	log := bytes.TrimSuffix(output, []byte{'\n'})
	if len(log) == 0 {
		return nil, nil
	}
	commits := bytes.Split(log, []byte{'\n'})
	return &commitIter{
		commits:    commits,
		repository: driver.repository,
	}, nil
}

func (driver *Git) RevList(options RevListOptions) (models.Iterator[models.Commit], error) {
	return nil, nil
}

func NewGit(options *config.Options) *Git {
	ctx := context.Background()
	return &Git{
		ctx:     ctx,
		options: options,
	}
}
