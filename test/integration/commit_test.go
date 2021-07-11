//go:build integration || commit

package main

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/models"
	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestCommitPattern(t *testing.T) {
	assert := assert.New(t)
	goldenFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/commit_message_pattern"))
	goldenFile = strings.Replace(goldenFile, "develop", Version, -1)
	defaultPattern := `(?m).*`
	conventionalCommitPattern := `(?m)^(build|ci|docs|feat|fix|perf|refactor|style|test)\([a-z]+\):\s([a-z\.\-\s]+)`
	multilineMessage := `
Merge branch 'feature/WATCHDOG-213' into 'develop'

feature/WATCHDOG-213

See merge request groupe-edf/watchdog#16`
	tests := []struct {
		name             string
		commitSubject    string
		pattern          string
		severity         models.Score
		rejectionMessage string
		skip             string
	}{
		{"Standard", "Initial commit", defaultPattern, models.SeverityLow, "", ""},
		{"CharacterEncoding", "Ajouter un pipeline d'intégration", defaultPattern, models.SeverityLow, "", ""},
		{"ConventionalCommit", "feat(scope): add new feature", conventionalCommitPattern, models.SeverityLow, "", ""},
		{"MultilineCommit", multilineMessage, defaultPattern, models.SeverityLow, "", ""},
		{"UnconventionalCommit", "This is SPARTA", conventionalCommitPattern, models.SeverityHigh, "Message must be formatted like type(scope): subject", ""},
		{"SkipPattern", "Merge branch 'feature/add-gitignore-file' into 'master'", conventionalCommitPattern, models.SeverityLow, "", "Merge branch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitHooksFile := fmt.Sprintf(goldenFile, test.pattern, test.rejectionMessage, test.skip)
			buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), test.commitSubject, nil)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity != models.SeverityLow {
				assert.Error(err)
				assert.Equal(1, len(issues))
				assert.Equal(test.severity, issues[0].Severity)
				assert.Equal(test.rejectionMessage, issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
			Suite.ResetLastCommit()
		})
	}
}

// TODO: Test all operators
func TestCommitLengthRule(t *testing.T) {
	assert := assert.New(t)
	goldenFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/commit_message_length"))
	goldenFile = strings.Replace(goldenFile, "develop", Version, -1)
	tests := []struct {
		name             string
		commitSubject    string
		operator         string
		length           string
		severity         models.Score
		rejectionMessage string
	}{
		{"EqualError", "Initial commit", "eq", "32", models.SeverityHigh, "Commit message not equal to {{ .Operand }}"},
		{"EqualSuccess", "Initial commit", "eq", "14", models.SeverityLow, ""},
		{"GreaterOrEqual", "Initial commit", "ge", "32", models.SeverityHigh, "Commit message longer or equal than {{ .Operand }}"},
		{"GreaterThan", "Initial commit", "gt", "32", models.SeverityHigh, "Commit message longer than {{ .Operand }}"},
		{"LowerOrEqual", "Initial commit", "le", "8", models.SeverityHigh, "Commit message shorter or equal than {{ .Operand }}"},
		{"LowerThan", "Initial commit", "lt", "8", models.SeverityHigh, "Commit message shorter than {{ .Operand }}"},
		{"NotEqualError", "Initial commit", "ne", "14", models.SeverityHigh, "Commit message equal to {{ .Operand }}"},
		{"NotEqualSuccess", "Initial commit", "ne", "32", models.SeverityLow, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitHooksFile := fmt.Sprintf(goldenFile, test.operator+" "+test.length, test.rejectionMessage)
			buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), test.commitSubject, nil)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity != models.SeverityLow {
				assert.Error(err)
				assert.Equal(ErrorPreReceiveHookDeclined, err)
				assert.Equal(1, len(issues))
				assert.Equal(test.severity, issues[0].Severity)
				var message bytes.Buffer
				tmpl := template.Must(template.New("").Parse(test.rejectionMessage))
				if err := tmpl.Execute(&message, issue.Data{Operand: test.length}); err != nil {
					t.Fatalf("Error redering message : %v", err)
				}
				assert.Equal(message.String(), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
			Suite.ResetLastCommit()
		})
	}
}

func TestCommitEmailRule(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/commit_email"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	signature := &object.Signature{
		Name:  "Habib MAALEM",
		Email: "habib.maalem@gmail.com",
		When:  time.Now(),
	}
	buffer, err := Suite.PushFile("master", ".githooks.yml", []byte(gitHooksFile), "Check author email in commit message", signature)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.Error(err)
	assert.Equal(ErrorPreReceiveHookDeclined, err)
	assert.Equal(1, len(issues))
	assert.Equal(models.SeverityHigh, issues[0].Severity)
	assert.Equal("Author email 'habib.maalem@gmail.com' is not valid email address", issues[0].Message)
	Suite.ResetLastCommit()
}
