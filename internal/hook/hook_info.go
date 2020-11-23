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
	Action   string
	NewRev   *object.Commit // New object name to be stored in the ref. When you delete a ref, this equals 40 zeroes.
	OldRev   *object.Commit // Old object name stored in the ref. When you create a new ref, this equals 40 zeroes.
	Ref      string         // The full name of the ref.
	RefType  string         // One of : heads / remotes / tags
	RefName  string
	RepoName string
	RepoPath string
}

// GetBranchName return the branche name
func (info *Info) GetBranchName() string {
	return info.RefName
}

// ParseHookAction return hook action
func ParseHookAction(info Info) string {
	action := "push"
	context := "branch"
	if info.RefType == "tags" {
		context = "tag"
	}
	if info.OldRev == nil && info.NewRev != nil {
		action = "create"
	} else if info.OldRev != nil && info.NewRev == nil {
		action = "delete"
	}
	return fmt.Sprintf("%s.%s", context, action)
}

// ParseInfo parse hook <old-value> SP <new-value> SP <ref-name> LF to models.Info
func ParseInfo(repository *git.Repository, input string) (*Info, error) {
	if input == "" {
		return nil, errors.New("Hook data is mandatory")
	}
	info, err := ReadHookInput(repository, strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	return info, nil
}

// ReadHookInput parse git hook input data and return Info object
func ReadHookInput(repository *git.Repository, input io.Reader) (*Info, error) {
	reader := bufio.NewReader(input)
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	chunks := strings.Split(strings.TrimSpace(string(line)), " ")
	if len(chunks) != 3 {
		return nil, fmt.Errorf("invalid hook input")
	}
	refchunks := strings.Split(chunks[2], "/")
	dir, _ := os.Getwd()
	info := Info{
		Ref:      chunks[2],
		RefType:  refchunks[1],
		RepoName: filepath.Base(dir),
		RepoPath: dir,
	}
	oldRevHash := plumbing.NewHash(chunks[0])
	if oldRevHash != plumbing.ZeroHash {
		commit, err := repository.CommitObject(oldRevHash)
		if err != nil {
			return nil, err
		}
		info.OldRev = commit
	}
	newRevHash, err := repository.ResolveRevision(plumbing.Revision(chunks[1]))
	if err != nil {
		return nil, err
	}
	if *newRevHash != plumbing.ZeroHash {
		commit, err := repository.CommitObject(*newRevHash)
		if err != nil {
			return nil, err
		}
		info.NewRev = commit
	}
	info.RefName = strings.Join(refchunks[2:], "/")
	info.Action = ParseHookAction(info)
	return &info, nil
}
