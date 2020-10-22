package util

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/groupe-edf/watchdog/internal/hook"
)

// FetchCommits fetching commits
func FetchCommits(repository *git.Repository, info *hook.Info, hookType string) (commits []*object.Commit, err error) {
	if hookType == "pre-receive" {
		return RevList(repository, info)
	}
	logOptions := &git.LogOptions{}
	if info != nil {
		if info.OldRev != nil {
			logOptions.Since = &info.OldRev.Committer.When
		}
		if info.NewRev != nil {
			logOptions.Until = &info.NewRev.Committer.When
		}
	}
	cmmitIter, err := repository.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	defer cmmitIter.Close()
	err = cmmitIter.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})
	return commits, err
}

// LoadRepository load git repository
func LoadRepository(uri string) (*git.Repository, error) {
	var repository *git.Repository
	var backend storage.Storer = memory.NewStorage()
	var err error
	if strings.Contains(uri, "://") || regexp.MustCompile(string(`^[A-Za-z]\w*@[A-Za-z0-9][\w.]*:`)).MatchString(uri) {
		cloneOptions := &git.CloneOptions{
			Progress: os.Stderr,
			URL:      uri,
		}
		repository, err = git.Clone(backend, nil, cloneOptions)
	} else {
		path, _ := os.Getwd()
		if uri[len(path)-1] == os.PathSeparator {
			path = path[:len(path)-1]
		}
		repository, err = git.PlainOpen(path)
	}
	return repository, err
}

// ParseGitPushOptions parse git push options
func ParseGitPushOptions() {
	// TODO: implement this
}

// RevListOptions defines the rules of rev-list func
type RevListOptions struct {
	// OldRev is the first reference hash to link
	OldRev plumbing.Hash
	// NewRev is the second reference hash to link
	NewRev plumbing.Hash
}

// RevList is native implemetation of git rev-list command
func RevList(repository *git.Repository, info *hook.Info) ([]*object.Commit, error) {
	var err error
	opts := RevListOptions{}
	if info.OldRev != nil {
		opts.OldRev = info.OldRev.Hash
	}
	if info.NewRev != nil {
		opts.NewRev = info.NewRev.Hash
	}
	commits := make([]*object.Commit, 0)
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
	for _, hash := range newRevObjects {
		commit, err := repository.CommitObject(hash)
		if err != nil {
			continue
		}
		commits = append(commits, commit)
	}
	return commits, err
}
