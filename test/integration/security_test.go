// +build integration security

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

func TestSecretRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		provider    string
		fileName    string
		fileContent string
		weakString  string
	}{
		{"AWS_ACCESS_KEY", "AWSProvider.java", `private static final String AWS_SECRET_ACCESS_KEY = "AKIAYYYYYYYYYYYYYYYY";`, `ACCESS_KEY = "AKIAYYYYYYYYYYYYYYYY"`},
		{"BASIC_AUTHENTICATION_URI", "application.json", `{"baseUrl": "https://www.acme.com", "assets": "https://root:Pa$$w0rd@acme.com:3306"}`, "https://root:Pa$$w0rd@acme.com:3306"},
		{"BASIC_AUTHENTICATION_HEADER", "main.py", `requests.get('https://www.acme.com', headers={'Authorization': 'Basic PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg=='})`, "Basic PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg=="},
		{"BASIC_AUTHENTICATION_HEADER", "deploy.sh", `set -e | curl -i -H 'Authorization: Basic PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg==' $NEXUS/watchdog/$VERSION/watchdog-$VERSION-linux_amd64.bin`, "Authorization: Basic PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg=="},
		{"CONFIDENTIAL", "README.md", "CONFIDENTIAL", "CONFIDENTIAL"},
		{"MYSQL", "database.json", `{"database": "mysql://root:Pa$$w0rd@localhost:3306"}`, "mysql://root:Pa$$w0rd@localhost:3306"},
		{"NPM_AUTHENTICATION", ".npmrc", "_auth = PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg==", "_auth = PFVTRVJOQU1FPjo8UEFTU1dPUkQ+Cg=="},
		{"PASSWORD", "application.properties", `jdbc:sqlserver://localhost;user=root;password=Pa$$w0rd;`, `password=Pa$$w0rd`},
		{"PASSWORD_EMPTY", ".env", `export DATABASE_PASSWORD=""`, ``},
		{"PASSWORD_QUOTES", "config.ini", `PASSWORD="Pa$$w0rd"`, `PASSWORD="Pa$$w0rd"`},
		{"PASSWORD_JSON", "config.json", `{"password": "Pa$$w0rd"}`, `"password": "Pa$$w0rd"`},
		{"PASSWORD_XML", "settings.xml", `<settings><servers><server><id>nexus</id><username>deployment</username><password>Pa$$w0rd@</password></server></servers></settings>`, `<password>Pa$$w0rd@</password>`},
		{"PRIVATE_KEY", "server.key", `-----BEGIN RSA PRIVATE KEY-----`, "-----BEGIN RSA PRIVATE KEY-----"},
		{"REDIS_URL", "database.json", `{"redisConnection": "redis://root:Pa$$w0rd@redis.acme.com:3306"}`, "redis://root:Pa$$w0rd@redis.acme.com:3306"},
		{"REDIS_URL_SSL", "database.json", `{"redisConnection": "rediss://root:Pa$$w0rd@redis.acme.com:3306"}`, "rediss://root:Pa$$w0rd@redis.acme.com:3306"},
		{"REDIS_URL_SOCKET", "database.json", `{"redisConnection": "redis-socket://root:Pa$$w0rd@redis.acme.com:3306"}`, "redis-socket://root:Pa$$w0rd@redis.acme.com:3306"},
		{"SECRET_KEY", ".env", `export API_TOKEN="1234567890abcdef"`, `API_TOKEN="1234567890abcdef"`},
		// FIXME: SECRET_KEY_VARIABLE should pass without errors
		{"SECRET_KEY_VARIABLE", ".env", `export accessKey="$secretKeyVariable"`, `accessKey="$secretKeyVariable"`},
	}
	gitHooksFile := helpers.LoadGolden(t, path.Join(RootDirectory, "/test/data/rules/security_secret.golden"))
	gitHooksFile = strings.Replace(gitHooksFile, "develop", Version, -1)
	for _, test := range tests {
		t.Run(test.provider, func(t *testing.T) {
			var files []helpers.File
			files = append(files, helpers.File{
				FileName:    ".githooks.yml",
				FileContent: []byte(gitHooksFile),
			}, helpers.File{
				FileName:    test.fileName,
				FileContent: []byte(test.fileContent),
			})
			buffer, err := Suite.CommitAndPush("master", files, "Add "+test.provider+" secret", nil)
			issues := helpers.ParseIssues(buffer.String())
			if test.weakString != "" {
				assert.Equal(ErrorPreReceiveHookDeclined, err)
				assert.Equal(1, len(issues))
				assert.Equal(issue.SeverityHigh, issues[0].Severity)
				assert.Equal(fmt.Sprintf("Secrets, token and passwords are forbidden, `%s:%s`", test.fileName, test.weakString), issues[0].Message)
			} else {
				assert.NoError(err)
				assert.Equal(0, len(issues))
			}
			Suite.ResetLastCommit()
		})
	}
}

func TestSecretRulesWithSkip(t *testing.T) {
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RulesDirectory, "security_secret_with_skip.golden"))
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
	issues := helpers.ParseIssues(buffer.String())
	assert.NoError(err)
	assert.Equal(0, len(issues))
}

func TestSkipWithGitPushOption(t *testing.T) {
	tearDownAll()
	setUpAll()
	assert := assert.New(t)
	gitHooksFile := helpers.LoadGolden(t, path.Join(RulesDirectory, "skip_with_git_push_option.golden"))
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
	issues := helpers.ParseIssues(buffer.String())
	assert.NoError(err)
	assert.Equal(0, len(issues))
}
