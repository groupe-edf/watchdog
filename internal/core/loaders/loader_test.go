package loaders

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/git"
)

func TestAPILoad(t *testing.T) {
	loader := NewAPILoader("http://localhost:3001", "GSHMG1A56JWNRX29YXE1IJQ0064QCXRL")
	policies, err := loader.LoadPolicies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 0 {
		t.Errorf("got %d policies, wanted %d policies", len(policies), 0)
	}
	rules, err := loader.LoadRules(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Errorf("got %d rules, wanted %d rules", len(rules), 0)
	}
}

func TestFileLoad(t *testing.T) {
	fileSystem := fstest.MapFS{
		".watchdog.yml": {Data: []byte("")},
	}
	loader := &FileLoader{
		fileSystem: fileSystem,
	}
	policies, err := loader.LoadPolicies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 0 {
		t.Errorf("got %d policies, wanted %d policies", len(policies), 0)
	}
	rules, err := loader.LoadRules(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Errorf("got %d rules, wanted %d rules", len(rules), 0)
	}
}

func TestGitFileLoad(t *testing.T) {
	backend := git.NewGit(&config.Options{})
	repository, _ := backend.Clone(context.Background(), git.CloneOptions{
		URL: "/home/maalem/workspaces/github/watchdog",
	})
	loader := &GitLoader{
		backend:    backend,
		repository: repository,
	}
	policies, err := loader.LoadPolicies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 0 {
		t.Errorf("got %d policies, wanted %d policies", len(policies), 0)
	}
	rules, err := loader.LoadRules(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Errorf("got %d rules, wanted %d rules", len(rules), 0)
	}
}
