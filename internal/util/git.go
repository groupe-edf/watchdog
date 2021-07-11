package util

import (
	"context"
	"fmt"
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
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/groupe-edf/watchdog/internal/backend"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/hook"
)

// FetchCommits fetching commits
func FetchCommits(repository *git.Repository, info *hook.Info, hookType string) (commitIter object.CommitIter, err error) {
	if hookType == "pre-receive" {
		return RevList(repository, info)
	}
	reference, _ := repository.Head()
	fmt.Printf("Running analysis on %s:", Colorize(Green, reference.Name().String()))
	logOptions := &git.LogOptions{}
	if info != nil {
		if info.OldRev != nil {
			logOptions.Since = &info.OldRev.Committer.When
		}
		if info.NewRev != nil {
			logOptions.Until = &info.NewRev.Committer.When
		}
	}
	commitIter, err = repository.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	return commitIter, err
}

// LoadRepository load git repository
func LoadRepository(ctx context.Context, options *config.Options) (*git.Repository, error) {
	var backend storage.Storer
	var err error
	uri := options.URI
	if strings.Contains(uri, "://") || regexp.MustCompile(string(`^[A-Za-z]\w*@[A-Za-z0-9][\w.]*:`)).MatchString(uri) {
		if options.CacheDirectory != "" {
			backend = filesystem.NewStorage(osfs.New(options.CacheDirectory), cache.NewObjectLRUDefault())
			_, err = os.Stat(options.CacheDirectory)
			if !os.IsNotExist(err) {
				os.RemoveAll(options.CacheDirectory)
			}
		} else {
			backend = memory.NewStorage()
		}
		cloneOptions := &git.CloneOptions{
			Progress: os.Stderr,
			URL:      uri,
		}
		if options.AuthBasicToken != "" {
			cloneOptions.Auth = &http.BasicAuth{
				Username: "watchdog",
				Password: options.AuthBasicToken,
			}
		}
		return git.CloneContext(ctx, backend, nil, cloneOptions)
	} else if stat, err := os.Stat(uri); err == nil && !stat.IsDir() {
		fs := osfs.New(filepath.Dir(uri))
		dot, _ := fs.Chroot(".git")
		storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
		return git.Open(storage, fs)
	} else {
		if uri[len(uri)-1] == os.PathSeparator {
			uri = uri[:len(uri)-1]
		}
		return git.PlainOpen(uri)
	}
}

// ParseGitPushOptions parse git push options
func ParseGitPushOptions() {
	// TODO: implement git push options
}

// RevList is native implemetation of git rev-list command
func RevList(repository *git.Repository, info *hook.Info) (object.CommitIter, error) {
	fmt.Printf("Running analysis on %s:", Colorize(Green, info.Ref.String()))
	var err error
	opts := backend.RevListOptions{}
	if info.OldRev != nil {
		opts.OldRev = info.OldRev.Hash
	}
	if info.NewRev != nil {
		opts.NewRev = info.NewRev.Hash
	}
	oldRevObjects := make([]plumbing.Hash, 0)
	newRevObjects := make([]plumbing.Hash, 0)
	if opts.OldRev != plumbing.ZeroHash {
		oldRevObjects, err = revlist.Objects(repository.Storer, []plumbing.Hash{opts.OldRev}, nil)
		if err != nil {
			return nil, err
		}
	}
	if opts.NewRev != plumbing.ZeroHash {
		newRevObjects, err = revlist.Objects(repository.Storer, []plumbing.Hash{opts.NewRev}, oldRevObjects)
		if err != nil {
			return nil, err
		}
	}
	commitIter := object.NewCommitIter(repository.Storer, storer.NewEncodedObjectLookupIter(repository.Storer, plumbing.AnyObject, newRevObjects))
	return commitIter, err
}
