package backend

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/groupe-edf/watchdog/internal/config"
)

type Backend struct {
	headReference plumbing.ReferenceName
	options       *config.Options
	repository    *git.Repository
	storer        storage.Storer
}

// RevListOptions defines the rules of rev-list func
type RevListOptions struct {
	// OldRev is the first reference hash to link
	OldRev plumbing.Hash
	// NewRev is the second reference hash to link
	NewRev plumbing.Hash
}

func (backend *Backend) Clone(ctx context.Context, cloneOptions *git.CloneOptions) error {
	if cloneOptions.URL == "" {
		return errors.New("url is mandatory")
	}
	var cachePath string
	var repository *git.Repository
	headReference, err := GetHead(cloneOptions)
	if err != nil {
		return err
	}
	backend.headReference = headReference
	cloneOptions.Progress = io.Discard
	if strings.Contains(cloneOptions.URL, "://") || regexp.MustCompile(string(`^[A-Za-z]\w*@[A-Za-z0-9][\w.]*:`)).MatchString(cloneOptions.URL) {
		if backend.options.CacheDirectory != "" {
			hasher := sha1.New()
			hasher.Write([]byte(cloneOptions.URL))
			cachePath = filepath.Join(backend.options.CacheDirectory, hex.EncodeToString(hasher.Sum(nil)))
			backend.storer = filesystem.NewStorage(osfs.New(cachePath), cache.NewObjectLRUDefault())
			_, err = os.Stat(cachePath)
			if !os.IsNotExist(err) {
				os.RemoveAll(cachePath)
			}
		} else {
			backend.storer = memory.NewStorage()
		}
		repository, err = git.CloneContext(ctx, backend.storer, nil, cloneOptions)
		if err != nil {
			return err
		}
	} else {
		repository, err = git.PlainOpen(cloneOptions.URL)
		if err != nil {
			return err
		}
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

func (backend *Backend) Head() (*plumbing.Reference, error) {
	return backend.repository.Head()
}

func (backend *Backend) RevList(options RevListOptions) (commitIter object.CommitIter, err error) {
	oldRevObjects := make([]plumbing.Hash, 0)
	newRevObjects := make([]plumbing.Hash, 0)
	if !options.OldRev.IsZero() {
		oldRevObjects, err = revlist.Objects(backend.repository.Storer, []plumbing.Hash{options.OldRev}, nil)
		if err != nil {
			return nil, err
		}
	}
	if !options.NewRev.IsZero() {
		newRevObjects, err = revlist.Objects(backend.repository.Storer, []plumbing.Hash{options.NewRev}, oldRevObjects)
		if err != nil {
			return nil, err
		}
	}
	commitIter = object.NewCommitIter(
		backend.repository.Storer,
		storer.NewEncodedObjectLookupIter(backend.repository.Storer, plumbing.AnyObject, newRevObjects),
	)
	return commitIter, err
}

func GetHead(options *git.CloneOptions) (plumbing.ReferenceName, error) {
	endpoint, err := transport.NewEndpoint(options.URL)
	if err != nil {
		return "", err
	}
	gitClient, err := client.NewClient(endpoint)
	if err != nil {
		return "", err
	}
	session, err := gitClient.NewUploadPackSession(endpoint, options.Auth)
	if err != nil {
		return "", err
	}
	info, err := session.AdvertisedReferences()
	if err != nil {
		return "", err
	}
	references, err := info.AllReferences()
	if err != nil {
		return "", err
	}
	return references["HEAD"].Target(), nil
}

func New(options *config.Options) *Backend {
	return &Backend{
		options: options,
	}
}
