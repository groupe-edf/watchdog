package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/groupe-edf/watchdog/internal/issue"
)

// CreateDummyFile create a dummy file
func CreateDummyFile(t *testing.T, size int64) []byte {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := file.Truncate(size); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, size)
	count, err := file.Read(data)
	if err != nil {
		t.Fatal(err)
	}
	return data[:count]
}

// LoadGolden load golden file
func LoadGolden(t *testing.T, goldenFile string) string {
	content, err := ioutil.ReadFile(filepath.Clean(goldenFile))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

// LoadInput load input file
func LoadInput(t *testing.T, inputFile string) []string {
	content, err := ioutil.ReadFile(filepath.Clean(inputFile))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n")
}

// ParseIssues parse integration test output
func ParseIssues(output string) []issue.Issue {
	issues := make([]issue.Issue, 0)
	slices := regexp.MustCompile(`(?s)([\-]{5}[A-Z ]+[\-]{5})(.+)([\-]{5}[A-Z ]+[\-]{5})`).FindStringSubmatch(output)
	if len(slices) > 0 {
		lines := regexp.MustCompile("\n+").Split(strings.TrimSpace(slices[2]), -1)
		for _, line := range lines {
			matches := regexp.MustCompile(`(?m)([a-z]+)=(?:(?:"(.*)")|(?:(?:([^\s]+)[\s])))`).FindAllStringSubmatch(line, -1)
			if len(matches) == 5 {
				issues = append(issues, issue.Issue{
					Hash:     matches[3][3],
					Message:  matches[4][2],
					Severity: issue.ParseScore(matches[0][3]),
				})
			}
		}
	}
	return issues
}
