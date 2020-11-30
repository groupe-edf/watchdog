// +build integration branch cli commit file jira security tag

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

const (
	OutputFormat = "json"
)

var (
	ErrorPreReceiveHookDeclined = errors.New("command error on refs/heads/master: pre-receive hook declined")
	RootDirectory               string
	RulesDirectory              string
	Suite                       *helpers.GitSuite
	Version                     string
)

func TestMain(m *testing.M) {
	setUpAll()
	exitStatus := m.Run()
	tearDownAll()
	os.Exit(exitStatus)
}

func TestNoGitHooksFile(t *testing.T) {
	tearDownAll()
	setUpAll()
	assert := assert.New(t)
	buffer, err := Suite.PushFile("master", "README.md", []byte("#Test Application"), "Commit with no .githooks file", nil)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.NoError(err)
	assert.Equal(0, len(issues))
}

func TestDefaultRules(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/default_empty_hooks"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), "Add .githooks.yml with no hooks", nil)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.NoError(err)
	assert.Equal(0, len(issues))
}

func setUpAll() {
	_, directory, _, _ := runtime.Caller(0)
	RootDirectory = filepath.Join(filepath.Dir(directory), "../..")
	err := os.Chdir(RootDirectory)
	if err != nil {
		fmt.Printf("Could not change directory %v", err)
		os.Exit(1)
	}
	versionFile, _ := ioutil.ReadFile(filepath.Join(RootDirectory, "VERSION"))
	Version = strings.TrimSpace(string(versionFile))
	RulesDirectory = RootDirectory + "/test/data/rules"
	make := exec.Command("make", "build-test")
	err = make.Run()
	if err != nil {
		fmt.Printf("Could not build watchdog CLI %v", err)
		os.Exit(1)
	}
	Suite = &helpers.GitSuite{
		RootDirectory: RootDirectory,
		OutputFormat:  OutputFormat,
	}
	err = Suite.SetUp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// setUp run before each test case
func setUp(t *testing.T) {

}

func tearDownAll() {
	err := Suite.Clean()
	if err != nil {
		fmt.Printf("Something went wrong when cleaning test files %v", err)
	}
}

// tearDown run after each test case
func tearDown(t *testing.T) {
	// TODO: revert commit after each test, not supported yel in go-git
}
