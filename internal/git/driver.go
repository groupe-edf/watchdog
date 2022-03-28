package git

import (
	"context"
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

var (
	ErrFailedToParseURL = errors.New("failed to parse url")
	ErrSchemeNotValid   = errors.New("scheme is not a valid transport")
)

type Protocols struct {
	protocols []string
}

func (protocols *Protocols) Valid(scheme string) bool {
	for _, protocol := range protocols.protocols {
		if protocol == scheme {
			return true
		}
	}
	return false
}

type CloneOptions struct {
	Authentication *BasicAuthentication
	Bare           bool
	Branch         string
	Depth          int
	Mirror         bool
	Quiet          bool
	Timeout        time.Duration
	URL            string
}

type BasicAuthentication struct {
	Username string
	Password string
}

type LogOptions struct {
}

type RevListOptions struct {
	OldRev string
	NewRev string
}

type Driver interface {
	Clone(ctx context.Context, cloneOptions CloneOptions) (*Repository, error)
	Commits(ctx context.Context, options LogOptions) (models.Iterator[models.Commit], error)
	File(filePath string) (string, error)
	GetRepository() *Repository
	Head() (string, error)
	RevList(options RevListOptions) (models.Iterator[models.Commit], error)
}

type URL struct {
	protocol string
}

func ParseURL(rawURL string) (*url.URL, error) {
	protocols := &Protocols{
		protocols: []string{
			"ssh",
			"git",
			"git+ssh",
			"http",
			"https",
			"ftp",
			"ftps",
			"rsync",
			"file",
		},
	}
	u, err := url.Parse(rawURL)
	if err == nil && protocols.Valid(u.Scheme) {
		err = ErrSchemeNotValid
	}
	return u, err
}

type commitIter struct {
	commits    [][]byte
	index      int
	repository *Repository
}

func (iterator *commitIter) Close() {
}

func (iterator *commitIter) ForEach(callback func(*models.Commit) error) error {
	for {
		commit, err := iterator.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		commit.Repository = &models.Repository{
			Storage: iterator.repository.Path(),
		}
		if err := callback(commit); err != nil {
			return err
		}
	}
}

func (iterator *commitIter) Next() (*models.Commit, error) {
	if iterator.index < len(iterator.commits) {
		commit, err := GetCommit(context.Background(), iterator.repository, string(iterator.commits[iterator.index]))
		if err != nil {
			return nil, err
		}
		iterator.index++
		return commit, nil
	}
	return nil, io.EOF
}
