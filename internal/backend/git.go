package backend

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Backend struct {
	repository *git.Repository
}

func (backend *Backend) Clone(ctx context.Context, uri string) error {
	var storer storage.Storer
	if strings.Contains(uri, "://") || regexp.MustCompile(string(`^[A-Za-z]\w*@[A-Za-z0-9][\w.]*:`)).MatchString(uri) {
		storer = memory.NewStorage()
	}
	cloneOptions := &git.CloneOptions{
		Progress: io.Discard,
		URL:      uri,
	}
	repository, err := git.CloneContext(ctx, storer, nil, cloneOptions)
	if err != nil {
		return err
	}
	backend.repository = repository
	return nil
}

func (backend *Backend) Commits(ctx context.Context, options *git.LogOptions) (object.CommitIter, error) {
	commitIter, err := backend.repository.Log(options)
	if err != nil {
		return nil, err
	}
	return commitIter, err
}

func New() *Backend {
	return &Backend{}
}
