// +build integration jira

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

func TestJiraRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name             string
		commitSubject    string
		conditionType    string
		skip             string
		severity         issue.Score
		rejectionMessage string
	}{
		{"SkipJiraIssueRule", "TECH Add .githooks.yml file", "issue", "TECH", issue.SeverityLow, ""},
		{"ValidCommitWithJiraIssue", "Add .githooks.yml file #WAT-1", "issue", "", issue.SeverityLow, ""},
		{"MissingJiraIssue", " Add .githooks.yml file", "issue", "", issue.SeverityHigh, "Commit message is missing the JIRA Issue 'JIRA-123'"},
	}
	goldenFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/jira_issue"))
	goldenFile = strings.Replace(goldenFile, "develop", Version, -1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitHooksFile := fmt.Sprintf(goldenFile, test.conditionType, test.skip, test.rejectionMessage)
			buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), test.commitSubject, nil)
			if err != nil {
				if test.severity == issue.SeverityHigh {
					assert.Equal(ErrorPreReceiveHookDeclined, err)
				}
				issues := helpers.ParseIssues(buffer.String(), OutputFormat)
				assert.Equal(test.severity, issues[0].Severity)
				assert.Equal(test.rejectionMessage, issues[0].Message)
			}
			Suite.ResetLastCommit()
		})
	}
}
