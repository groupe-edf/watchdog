//go:build integration || file
// +build integration file

package main

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/c2h5oh/datasize"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/issue"
	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestFileExist(t *testing.T) {
	// TODO: implement test
}

func TestFileExtensionNotAllowedRule(t *testing.T) {
	var files []helpers.File
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/file_extension_not_allowed"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	files = append(files, helpers.File{
		FileName:    ".githooks.yml",
		FileContent: []byte(gitHooksFile),
	})
	files = append(files, helpers.File{
		FileName:    "explorer.exe",
		FileContent: []byte(""),
	})
	buffer, err := Suite.CommitAndPush("master", files, "Add .githooks.yml to exclude extensions", nil)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.Error(err)
	assert.Equal(ErrorPreReceiveHookDeclined, err)
	assert.Equal(1, len(issues))
	assert.Equal(models.SeverityHigh, issues[0].Severity)
	assert.Equal("'*.exe' files are not allowed", issues[0].Message)
	Suite.ResetLastCommit()
}

func TestFileSizeExceededRule(t *testing.T) {
	assert := assert.New(t)
	var files []helpers.File
	goldenFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/file_size_exceeded"))
	goldenFile = strings.Replace(goldenFile, "develop", Version, -1)
	tests := []struct {
		name             string
		operator         string
		size             string
		fileSize         string
		severity         models.Score
		rejectionMessage string
	}{
		{"LowerThanSuccess", "lt", "5mb", "2mb", models.SeverityLow, ""},
		{"LowerThanError", "lt", "5mb", "10mb", models.SeverityHigh, "File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitHooksFile := fmt.Sprintf(goldenFile, test.operator+" "+test.size, test.rejectionMessage)
			files = append(files, helpers.File{
				FileName:    ".githooks.yml",
				FileContent: []byte(gitHooksFile),
			})
			var fileSize datasize.ByteSize
			err := fileSize.UnmarshalText([]byte(test.fileSize))
			if err != nil {
				t.Fatal("Error parsing file size")
			}
			files = append(files, helpers.File{
				FileName:    "postgresql.jar",
				FileContent: helpers.CreateDummyFile(t, int64(fileSize.Bytes())),
			})
			var size datasize.ByteSize
			err = size.UnmarshalText([]byte(test.size))
			if err != nil {
				t.Fatal("Error parsing file size")
			}
			buffer, err := Suite.CommitAndPush("master", files, "Add database dependency driver postgresql.jar", nil)
			issues := helpers.ParseIssues(buffer.String(), OutputFormat)
			if test.severity != models.SeverityLow {
				assert.Error(err)
				assert.Equal(ErrorPreReceiveHookDeclined, err)
				assert.Equal(1, len(issues))
				assert.Equal(models.SeverityHigh, issues[0].Severity)
				var message bytes.Buffer
				tmpl := template.Must(template.New("").Parse(test.rejectionMessage))
				if err := tmpl.Execute(&message, issue.Data{
					Object:  "postgresql.jar",
					Operand: size.HumanReadable(),
					Value:   fileSize.HumanReadable(),
				}); err != nil {
					t.Fatalf("Error redering message : %v", err)
				}
				assert.Equal(message.String(), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
		})
		Suite.ResetLastCommit()
	}
}
