// +build tag

package main

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/groupe-edf/watchdog/internal/issue"
	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTagRules(t *testing.T) {
	var files []helpers.File
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/tag_semver"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	files = append(files, helpers.File{
		FileName:    ".githooks.yml",
		FileContent: []byte(gitHooksFile),
	})
	_, err := Suite.CommitAndPush("master", files, "Add .githooks.yml to check tags", nil)
	if err != nil {
		t.Fatalf("Something went wrong when committing .giithooks.yml file %v", err)
	}
	tests := []struct {
		name     string
		severity issue.Score
	}{
		{"release", issue.SeverityHigh},
		{"2", issue.SeverityHigh},
		{"2.0", issue.SeverityHigh},
		{"2.0.0", issue.SeverityLow},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buffer, err := Suite.AnnotatedTag(test.name, Suite.LastCommit)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity == issue.SeverityHigh {
				fmt.Print(buffer.String())
				assert.Equal(fmt.Errorf("command error on refs/tags/%s: pre-receive hook declined", test.name), err)
				assert.Equal(1, len(issues))
				assert.Equal(test.severity, issues[0].Severity)
				assert.Equal(fmt.Sprintf("Tag version `%s` must respect semantic versionning v2.0.0 https://semver.org/", test.name), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
		})
	}
}
