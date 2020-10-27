package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitSuite data structure
type GitSuite struct {
	BarePath       string
	ClonePath      string
	LastCommit     plumbing.Hash
	OutputFormat   string
	PreviousCommit plumbing.Hash
	Repository     *git.Repository
	RootDirectory  string
	Server         *git.Repository
	tempDirectory  string
}

// File data structure
type File struct {
	FileName    string
	FileContent []byte
}

var (
	defaultSignature = object.Signature{
		Name:  "Habib MAALEM",
		Email: "habib.maalem@gmail.com",
		When:  time.Now(),
	}
)

// Clean clean data
func (suite *GitSuite) Clean() error {
	err := os.RemoveAll(suite.BarePath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(suite.ClonePath)
	if err != nil {
		return err
	}
	return nil
}

// CreateBranch create new branch
func (suite *GitSuite) CreateBranch(ref string) (*bytes.Buffer, error) {
	hash, err := suite.Repository.Head()
	if err != nil {
		return nil, err
	}
	// Create new ref
	name := plumbing.ReferenceName("refs/heads/" + ref)
	reference := plumbing.NewHashReference(name, hash.Hash())
	err = suite.Repository.Storer.SetReference(reference)
	if err != nil {
		return nil, err
	}
	return suite.Push(ref)
}

// LightweightTag create new lightweight tag
func (suite *GitSuite) LightweightTag(tag string, hash plumbing.Hash) (*bytes.Buffer, error) {
	name := plumbing.ReferenceName("refs/tags/" + tag)
	reference := plumbing.NewHashReference(name, hash)
	err := suite.Repository.Storer.SetReference(reference)
	if err != nil {
		return nil, err
	}
	refSpec := config.RefSpec("refs/tags/*:refs/tags/*")
	buffer := bytes.NewBuffer(nil)
	err = suite.Repository.Push(&git.PushOptions{
		RefSpecs:   []config.RefSpec{refSpec},
		RemoteName: "origin",
		Progress:   buffer,
	})
	return buffer, err
}

// AnnotatedTag create annotated tag
func (suite *GitSuite) AnnotatedTag(tag string, hash plumbing.Hash) (*bytes.Buffer, error) {
	tagObject := object.Tag{
		Name:         tag,
		Message:      "Release of " + tag,
		Tagger:       defaultSignature,
		PGPSignature: "",
		Target:       hash,
		TargetType:   plumbing.CommitObject,
	}
	object := suite.Repository.Storer.NewEncodedObject()
	err := tagObject.Encode(object)
	if err != nil {
		return nil, err
	}
	hash, err = suite.Repository.Storer.SetEncodedObject(object)
	if err != nil {
		return nil, err
	}
	err = suite.Repository.Storer.SetReference(plumbing.NewReferenceFromStrings("refs/tags/"+tag, hash.String()))
	if err != nil {
		return nil, err
	}
	refSpec := config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag))
	buffer := bytes.NewBuffer(nil)
	err = suite.Repository.Push(&git.PushOptions{
		RefSpecs:   []config.RefSpec{refSpec},
		RemoteName: "origin",
		Progress:   buffer,
	})
	return buffer, err
}

// Push to remote repository
func (suite *GitSuite) Push(ref string) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(nil)
	refSpec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", ref, ref))
	err := suite.Repository.Push(&git.PushOptions{
		RefSpecs:   []config.RefSpec{refSpec},
		RemoteName: "origin",
		Progress:   buffer,
	})
	return buffer, err
}

// CommitAndPush commit and push multiple files to bare repository
func (suite *GitSuite) CommitAndPush(ref string, files []File, commitMessage string, signature *object.Signature) (*bytes.Buffer, error) {
	commitHash, err := suite.addFile(commitMessage, signature, files)
	if err != nil {
		return nil, err
	}
	suite.LastCommit = commitHash
	return suite.Push(ref)
}

// PushFile commit and push single file to bare repository
func (suite *GitSuite) PushFile(ref string, fileName string, fileContent []byte, commitMessage string, signature *object.Signature) (*bytes.Buffer, error) {
	var files []File
	files = append(files, File{
		FileName:    fileName,
		FileContent: fileContent,
	})
	buffer, err := suite.CommitAndPush(ref, files, commitMessage, signature)
	return buffer, err
}

// ResetLastCommit git reset --hard [origin/master|hash
func (suite *GitSuite) ResetLastCommit() error {
	remoteRef, err := suite.Repository.Reference(plumbing.ReferenceName("refs/remotes/origin/master"), true)
	if err != nil {
		return err
	}
	worktree, err := suite.Repository.Worktree()
	if err != nil {
		return err
	}
	if err = worktree.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: remoteRef.Hash(),
	}); err != nil {
		return err
	}
	_, err = suite.Repository.Reference(plumbing.ReferenceName("HEAD"), false)
	if err != nil {
		return err
	}
	return nil
}

// SetUp initialize bare and clone git repository
func (suite *GitSuite) SetUp() error {
	suite.tempDirectory = filepath.Join(suite.RootDirectory, "/target")
	err := suite.SetUpBareRepository()
	if err != nil {
		return err
	}
	err = suite.SetUpCloneRepository()
	if err != nil {
		return err
	}
	return nil
}

// SetUpBareRepository create new bare repository
func (suite *GitSuite) SetUpBareRepository() error {
	barePath, err := ioutil.TempDir("", "repository.git")
	if err != nil {
		return err
	}
	suite.BarePath = barePath
	repository, err := git.PlainInit(barePath, true)
	if err != nil {
		return err
	}
	suite.Server = repository
	suite.installPreReceiveHook()
	return nil
}

// SetUpCloneRepository clone git repository
func (suite *GitSuite) SetUpCloneRepository() error {
	clonePath, err := ioutil.TempDir("", "repository")
	if err != nil {
		return err
	}
	suite.ClonePath = clonePath
	repository, err := git.PlainInit(clonePath, false)
	if err != nil {
		return err
	}
	suite.Repository = repository
	_, err = suite.Repository.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{suite.BarePath},
	})
	if err != nil {
		return err
	}
	return nil
}

func (suite *GitSuite) addFile(commitMessage string, signature *object.Signature, files []File) (plumbing.Hash, error) {
	worktree, err := suite.Repository.Worktree()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	for _, file := range files {
		err = ioutil.WriteFile(filepath.Join(suite.ClonePath, file.FileName), file.FileContent, 0777)
		if err != nil {
			return plumbing.ZeroHash, err
		}
		_, err = worktree.Add(file.FileName)
		if err != nil {
			return plumbing.ZeroHash, err
		}
	}
	if signature == nil {
		signature = &defaultSignature
	}
	// Create commit object
	commitHash, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: signature,
	})
	if err != nil {
		return plumbing.ZeroHash, err
	}
	// Process commit
	_, err = suite.Repository.CommitObject(commitHash)
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return commitHash, nil
}

// InstallPreReceiveHook install pre-receive hook file in hooks directory
func (suite *GitSuite) installPreReceiveHook() {
	hooks := filepath.Join(suite.BarePath, "hooks")
	err := os.MkdirAll(hooks, 0750)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	template, _ := ioutil.ReadFile(path.Join(suite.RootDirectory, "/test/data/pre-receive"))
	preReceiveHook := []byte(fmt.Sprintf(string(template), path.Join(suite.RootDirectory, "/target/bin/watchdog"), suite.OutputFormat))
	err = ioutil.WriteFile(filepath.Join(hooks, "pre-receive"), preReceiveHook, 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
