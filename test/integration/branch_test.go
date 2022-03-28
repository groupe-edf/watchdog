//go:build integration || branch
// +build integration branch

package main

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/groupe-edf/watchdog/internal/core/models"
	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestBranchNaming(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/branch_naming"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), "Add .githooks.yml with branch naming rules", nil)
	if err != nil {
		fmt.Print(buffer.String())
		t.Fatalf("Something went wrong when committing .giithooks.yml file %v", err)
	}
	tests := []struct {
		name     string
		severity models.Score
	}{
		{"documentation", models.SeverityHigh},
		{"production", models.SeverityHigh},
		{"hotfix", models.SeverityHigh},
		{"feature/update-documentation", models.SeverityLow},
		{"release/milestone-1.0.0", models.SeverityLow},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer, err := Suite.CreateBranch(test.name)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity == models.SeverityHigh {
				assert.Error(err)
				assert.Equal(fmt.Errorf("command error on refs/heads/%s: pre-receive hook declined", test.name), err)
				assert.Equal(1, len(issues))
				assert.Equal(fmt.Sprintf("Branch `%s` must match Gitflow naming convention", test.name), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
		})
	}
}

func TestBranchProtected(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/branch_protected"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), "Add .githooks.yml to check protected branches", nil)
	if err != nil {
		t.Fatalf("Something went wrong when committing .giithooks.yml file %v", err)
	}
	tests := []struct {
		name     string
		severity models.Score
	}{
		{"develop", models.SeverityLow},
		{"qa", models.SeverityLow},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer, err = Suite.CreateBranch(test.name)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity == models.SeverityHigh {
				assert.Equal(fmt.Errorf("command error on refs/heads/%s: pre-receive hook declined", test.name), err)
				assert.Equal(1, len(issues))
				assert.Equal(fmt.Sprintf("Branch %s is protected", test.name), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
		})
	}
}
