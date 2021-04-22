package hook

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	// BranchPushAction a branch is pushed
	BranchPushAction = "branch.push"
	// BranchCreateAction a branch is created
	BranchCreateAction = "branch.create"
	// BranchDeleteAction a branch is deleted
	BranchDeleteAction = "branch.delete"
	// TagCreateAction a tag is created
	TagCreateAction = "tag.create"
	// TagDeleteAction a tag is deleted
	TagDeleteAction = "tag.delete"
	// ZeroCommit hash if no commit
	ZeroCommit = "0000000000000000000000000000000000000000"
)

var (
	// ErrNoHookData no hook data error
	ErrNoHookData      = errors.New("Hook data is mandatory")
	ErrInvalidHookData = errors.New("invalid hook input")
	// HookTypes that are supported by Watchdog
	HookTypes = [...]string{
		"pre-receive",
		"update",
		"post-receive",
	}
)

// Hooks server side git hooks
type Hooks struct {
	PreReceive  string
	Update      string
	PostReceive string
}

// Info git hook data structure
type Info struct {
	Action     string
	NewRev     *object.Commit         // New object name to be stored in the ref. When you delete a ref, this equals 40 zeroes.
	OldRev     *object.Commit         // Old object name stored in the ref. When you create a new ref, this equals 40 zeroes.
	Ref        plumbing.ReferenceName // The full name of the ref.
	RefType    string                 // One of : heads / remotes / tags
	RepoName   string
	RepoPath   string
	repository *git.Repository
}

// ParseHookInput parse git hook input data and return Info object
// format: <old-value> SP <new-value> SP <ref-name> LF to models.Info
func (info *Info) ParseHookInput(input io.Reader) error {
	reader := bufio.NewReader(input)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	chunks := strings.Split(strings.TrimSpace(string(line)), " ")
	if len(chunks) != 3 {
		return ErrInvalidHookData
	}
	info.Ref = plumbing.ReferenceName(chunks[2])
	refChunks := strings.Split(chunks[2], "/")
	info.RefType = refChunks[1]
	oldRevHash := plumbing.NewHash(chunks[0])
	if oldRevHash != plumbing.ZeroHash {
		commit, err := info.repository.CommitObject(oldRevHash)
		if err != nil {
			return err
		}
		info.OldRev = commit
	}
	newRevHash, err := info.repository.ResolveRevision(plumbing.Revision(chunks[1]))
	if err != nil {
		return err
	}
	if *newRevHash != plumbing.ZeroHash {
		commit, err := info.repository.CommitObject(*newRevHash)
		if err != nil {
			return err
		}
		info.NewRev = commit
	}
	info.parseHookAction()
	return nil
}

// ParseHookAction return hook action
func (info *Info) parseHookAction() {
	action := "push"
	context := "branch"
	if info.Ref.IsTag() {
		context = "tag"
	}
	if info.OldRev == nil && info.NewRev != nil {
		action = "create"
	} else if info.OldRev != nil && info.NewRev == nil {
		action = "delete"
	}
	info.Action = fmt.Sprintf("%s.%s", context, action)
}

func ParseInfo(repository *git.Repository) (*Info, error) {
	directory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	info := &Info{
		Action:     "branch.push",
		RepoName:   filepath.Base(directory),
		RepoPath:   directory,
		repository: repository,
	}
	return info, nil
}
