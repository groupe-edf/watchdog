// +build integration security

package main

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/groupe-edf/watchdog/internal/issue"
	helpers "github.com/groupe-edf/watchdog/internal/test"
	"github.com/stretchr/testify/assert"
)

type Issue struct {
	Offender string
	Rule     string
}

func TestSecretRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Provider    string
		FileName    string
		FileContent string
		Count       int
		Issue       *Issue
	}{
		{"INITIAL_COMMIT", "README.md", "# Readme", 0, nil},
		{"BASE_64_AUTHORIZATION_HEADER", "deploy.sh", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/deploy.sh")), 1, &Issue{Offender: "PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg==", Rule: "BASE_64"}},
		{"BASE_64_JSON", "config.json", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/config.json")), 1, &Issue{Offender: "X3Rva2VuOjEyMzQ1Njc4OTBBQkNERUY=", Rule: "BASE_64"}},
		{"BASE_64_NPM", ".npmrc", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/.npmrc")), 3, &Issue{Offender: "PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg==", Rule: "BASE_64"}},
		{"CONFIDENTIAL", "SECURITY.md", "CONFIDENTIAL", 1, &Issue{Offender: "CONFIDENTIAL", Rule: "CONFIDENTIAL"}},
		{"CONNECTION_STRING", "application.properties", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/database.properties")), 7, &Issue{Offender: "Pa$$w0rd", Rule: "CONNECTION_STRING"}},
		{"CONNECTION_STRING_PIP", "pip.conf", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/pip.conf")), 2, &Issue{Offender: "Pa$$w0rd", Rule: "CONNECTION_STRING"}},
		{"ENTROPY_MYSQL_DUMP", "dump.sql", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/dump.sql")), 8, &Issue{Offender: "$2y$12$s3fn56ajUsYzNCVLkfprB.zHmMmOOBJ/Ro/wU0wRiIWaIRrk9gcei", Rule: "ENTROPY"}},
		{"HTPASSWD", ".htpasswd", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/.htpasswd")), 5, &Issue{Offender: "Pa$$w0rd", Rule: "HTPASSWD"}},
		{"LANGUAGE_GO", "main.go", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/language.go")), 1, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"LANGUAGE_JAVA", "AWSProvider.java", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/language.java")), 5, &Issue{Offender: "AKIAYYYYYYYYYYYYYYYY", Rule: "AWS_ACCESS_KEY"}},
		{"LANGUAGE_SHELL", "deploy.sh", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/language.sh")), 7, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"PASSWORD", "application.properties", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/application.properties")), 3, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"PASSWORD_ENV", ".env", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/.env")), 2, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"PASSWORD_JSON", "application.json", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/application.json")), 2, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"PASSWORD_QUOTES", "config.ini", `PASSWORD="Pa$$w0rd"`, 1, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD"}},
		{"PASSWORD_XML", "settings.xml", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/settings.xml")), 1, &Issue{Offender: "Pa$$w0rd", Rule: "PASSWORD_XML"}},
		{"PRIVATE_KEY", "server.key", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/server.key")), 1, &Issue{Offender: "-----BEGIN RSA PRIVATE KEY-----", Rule: "ASYMMETRIC_PRIVATE_KEY"}},
		{"SECRET_KEY_TOML", "config.toml", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/config.toml")), 2, &Issue{Offender: "1234567890ABCDEF", Rule: "SECRET_KEY"}},
		{"SECRET_KEY_ENV", ".env", helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/production.env")), 1, &Issue{Offender: "t&XG+FJ%M@8XYv5a!xaR", Rule: "SECRET_KEY"}},
	}
	rejectionMessage := "Secrets, token and passwords are forbidden, `{{ .Object }}:{{ Hide .Value 4 }}`"
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/security_secret"))
	gitHooksFile = fmt.Sprintf(gitHooksFile, Version)
	for _, test := range tests {
		t.Run(test.Provider, func(t *testing.T) {
			var files []helpers.File
			var issues = make([]issue.Issue, 0)
			files = append(files, helpers.File{
				FileName:    ".githooks.yml",
				FileContent: []byte(gitHooksFile),
			}, helpers.File{
				FileName:    test.FileName,
				FileContent: []byte(test.FileContent),
			})
			buffer, err := Suite.CommitAndPush("master", files, "Add "+test.Provider+" secret", nil)
			issues = helpers.ParseIssues(buffer.String(), OutputFormat)
			assert.Equal(test.Count, len(issues))
			if test.Issue != nil {
				assert.Equal(ErrorPreReceiveHookDeclined, err)
				assert.NotEmpty(issues[0].Leaks)
				assert.Equal(test.Issue.Rule, issues[0].Leaks[0].Rule)
				assert.Equal(issue.SeverityHigh, issues[0].Severity)
				var message bytes.Buffer
				t := template.Must(template.New("").Funcs(issue.FunctionsMap).Parse(rejectionMessage))
				_ = t.Execute(&message, issue.Data{
					Object: test.FileName,
					Value:  test.Issue.Offender,
				})
				assert.Equal(message.String(), issues[0].Message)
			} else {
				assert.NoError(err)

			}
			err = Suite.ResetLastCommit()
			if err != nil {
				t.Fatalf("Something went wrong when trying to reset last commit: %s", err)
			}
		})
	}
}

func TestSecretRulesWithSkip(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RulesDirectory, "security_secret_with_skip"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	var files []helpers.File
	files = append(files, helpers.File{
		FileName:    ".githooks.yml",
		FileContent: []byte(gitHooksFile),
	}, helpers.File{
		FileName:    "application.json",
		FileContent: []byte(`{"redisConnection": "rediss://root:Pa$$w0rd@redis.acme.com:3306"}`),
	})
	buffer, err := Suite.CommitAndPush("master", files, "Add REDIS_URL secret", nil)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.NoError(err)
	assert.Equal(0, len(issues))
}

func TestSkipWithGitPushOption(t *testing.T) {
	tearDownAll()
	setUpAll()
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RulesDirectory, "skip_with_git_push_option"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	var files []helpers.File
	files = append(files, helpers.File{
		FileName:    ".githooks.yml",
		FileContent: []byte(gitHooksFile),
	}, helpers.File{
		FileName:    "application.json",
		FileContent: []byte(`{"password": "Pa$$w0rd"}`),
	})
	buffer, err := Suite.CommitAndPush("master", files, "Add REDIS_URL secret [skip hooks.security.secret]", nil)
	issues := helpers.ParseIssues(buffer.String(), OutputFormat)
	assert.NoError(err)
	assert.Equal(0, len(issues))
}
