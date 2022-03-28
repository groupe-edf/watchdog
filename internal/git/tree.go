package git

import (
	"context"
	"fmt"
	"io"
	"math"
	"path"
	"strings"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

const (
	// EntryModeBlob
	EntryModeBlob EntryMode = 0o100644
	// EntryModeExec
	EntryModeExec EntryMode = 0o100755
	// EntryModeSymlink
	EntryModeSymlink EntryMode = 0o120000
	// EntryModeCommit
	EntryModeCommit EntryMode = 0o160000
	// EntryModeTree
	EntryModeTree EntryMode = 0o040000
)

type Tree struct {
	ID                     string
	entries                Entries
	entriesParsed          bool
	entriesRecursive       Entries
	entriesRecursiveParsed bool
	ResolvedID             string
	parentTree             *Tree
	repository             *Repository
}

func (tree *Tree) GetBlobByPath(ctx context.Context, relativePath string) (*Blob, error) {
	entry, err := tree.GetTreeEntryByPath(ctx, relativePath)
	if err != nil {
		return nil, err
	}
	if !entry.IsDir() && !entry.IsSubModule() {
		return entry.Blob(), nil
	}
	return nil, fmt.Errorf("blob %s not found", relativePath)
}

func (tree *Tree) GetTreeEntryByPath(ctx context.Context, relativePath string) (*TreeEntry, error) {
	if len(relativePath) == 0 {
		return &TreeEntry{
			parentTree: tree,
			ID:         tree.ID,
			name:       "",
			fullName:   "",
			entryMode:  EntryModeTree,
		}, nil
	}
	relativePath = path.Clean(relativePath)
	parts := strings.Split(relativePath, "/")
	var err error
	t := tree
	for index, name := range parts {
		if index == len(parts)-1 {
			entries, err := t.ListEntries(ctx)
			if err != nil {
				return nil, err
			}
			for _, entry := range entries {
				if entry.Name() == name {
					return entry, nil
				}
			}
		} else {
			t, err = t.SubTree(ctx, name)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, ErrCommitNotFound
}

func (tree *Tree) ListEntries(ctx context.Context) (Entries, error) {
	if tree.entriesParsed {
		return tree.entries, nil
	}
	if tree.repository != nil {
		writer, reader, cancel := CatFileBatch(ctx, tree.repository.Path())
		defer cancel()
		// Read header
		_, _ = writer.Write([]byte(tree.ID + "\n"))
		_, typ, size, err := ReadBatchLine(reader)
		if err != nil {
			return nil, err
		}
		if typ == "commit" {
			treeID, err := ReadTreeID(reader, size)
			if err != nil && err != io.EOF {
				return nil, err
			}
			_, _ = writer.Write([]byte(treeID + "\n"))
			_, typ, size, err = ReadBatchLine(reader)
			if err != nil {
				return nil, err
			}
		}
		if typ == "tree" {
			tree.entries, err = catBatchParseTreeEntries(tree, reader, size)
			if err != nil {
				return nil, err
			}
			tree.entriesParsed = true
			return tree.entries, nil
		}
		for size > math.MaxInt32 {
			discarded, err := reader.Discard(math.MaxInt32)
			size -= int64(discarded)
			if err != nil {
				return nil, err
			}
		}
		for size > 0 {
			discarded, err := reader.Discard(int(size))
			size -= int64(discarded)
			if err != nil {
				return nil, err
			}
		}
	}
	output, _, _ := NewCommand(ctx, "ls-tree", "-l", tree.ID).RunStdBytes(&RunOptions{Dir: tree.repository.Path()})
	var err error
	tree.entries, err = parseTreeEntries(output, tree)
	return tree.entries, err
}

func (tree *Tree) ListEntriesRecursive(ctx context.Context) (Entries, error) {
	stdout, _, runErr := NewCommand(ctx, "ls-tree", "-t", "-l", "-r", tree.ID).RunStdBytes(&RunOptions{Dir: tree.repository.Path()})
	if runErr != nil {
		return nil, runErr
	}
	var err error
	tree.entriesRecursive, err = parseTreeEntries(stdout, tree)
	if err == nil {
		tree.entriesRecursiveParsed = true
	}
	return tree.entriesRecursive, err
}

func (tree *Tree) SubTree(ctx context.Context, relativePath string) (*Tree, error) {
	if len(relativePath) == 0 {
		return tree, nil
	}
	paths := strings.Split(relativePath, "/")
	var (
		err   error
		g     = tree
		p     = tree
		entry *TreeEntry
	)
	for _, name := range paths {
		entry, err = p.GetTreeEntryByPath(ctx, name)
		if err != nil {
			return nil, err
		}
		g, err = GetRepositoryTree(ctx, tree.repository, entry.ID)
		if err != nil {
			return nil, err
		}
		g.parentTree = p
		p = g
	}
	return g, nil
}

type Entries []*TreeEntry

type EntryMode int

type TreeEntry struct {
	ID         string
	parentTree *Tree
	entryMode  EntryMode
	name       string
	size       int64
	sized      bool
	fullName   string
}

func (tree *TreeEntry) Blob() *Blob {
	return &Blob{
		ID:         tree.ID,
		name:       tree.Name(),
		repository: tree.parentTree.repository,
		size:       tree.size,
	}
}

func (tree *TreeEntry) IsDir() bool {
	return tree.entryMode == EntryModeTree
}

func (tree *TreeEntry) IsSubModule() bool {
	return tree.entryMode == EntryModeCommit
}

func (entry *TreeEntry) Name() string {
	if entry.fullName != "" {
		return entry.fullName
	}
	return entry.name
}

func NewTree(repository *Repository, commitID string) *Tree {
	return &Tree{
		ID:         commitID,
		repository: repository,
	}
}

func GetTree(commit *models.Commit) *Tree {
	repository, _ := NewRepository(commit.Repository.Storage)
	return NewTree(repository, commit.Hash)
}

func GetRepositoryTree(ctx context.Context, repository *Repository, id string) (*Tree, error) {
	wr, reader, cancel := CatFileBatch(ctx, repository.Path())
	defer cancel()
	_, _ = wr.Write([]byte(id + "\n"))
	_, typ, size, err := ReadBatchLine(reader)
	if err != nil {
		return nil, err
	}
	switch typ {
	case "commit":
		commit, err := CommitFromReader(repository, id, io.LimitReader(reader, size))
		if err != nil {
			return nil, err
		}
		if _, err := reader.Discard(1); err != nil {
			return nil, err
		}
		commit.Tree.(*Tree).ResolvedID = commit.Hash
		return commit.Tree.(*Tree), nil
	case "tree":
		tree := NewTree(repository, id)
		tree.ResolvedID = id
		tree.entries, err = catBatchParseTreeEntries(tree, reader, size)
		if err != nil {
			return nil, err
		}
		tree.entriesParsed = true
		return tree, nil
	default:
		return nil, fmt.Errorf("commit %s doesn't exist", id)
	}
}
