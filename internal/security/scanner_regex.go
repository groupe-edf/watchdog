package security

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
)

// IsFalsePositiveFunc function to check false positive secrets
type IsFalsePositiveFunc func(file string, line string, secret string) int

var (
	falsePostiveHeuristics = []IsFalsePositiveFunc{
		// Check if is dynamic variable
		IsFalsePositive,
	}
)

// RegexScanner data struct
type RegexScanner struct {
	Logger    logging.Interface
	Rules     []models.Rule
	Whitelist models.Whitelist
}

// Scan scan commit and return list of leaks if found
func (scanner *RegexScanner) Scan(commit *object.Commit) (leaks []models.Leak, err error) {
	if len(commit.ParentHashes) == 0 {
		return
	}
	parent, err := commit.Parent(0)
	if err != nil {
		return
	}
	patch, err := parent.Patch(commit)
	if err != nil {
		return
	}
	patchContent := patch.String()
	for _, filePatch := range patch.FilePatches() {
		if filePatch.IsBinary() {
			continue
		}
		_, to := filePatch.Files()
		if len(scanner.Whitelist.Files) != 0 && to != nil {
			for _, fileRegex := range scanner.Whitelist.Files {
				if regexp.MustCompile(fileRegex).FindString(to.Path()) != "" {
					return
				}
			}
		}
		for _, chunk := range filePatch.Chunks() {
			if chunk.Type() != diff.Delete {
				lineLookup := make(map[string]bool)
				for _, rule := range scanner.Rules {
					offenders := rule.Scan(chunk.Content())
					if len(offenders) == 0 {
						continue
					}
					leaks = append(leaks, models.Leak{
						Author:     commit.Author.Name,
						CommitHash: commit.Hash.String(),
						CreatedAt:  commit.Author.When,
						File:       to.Path(),
						Line:       offenders[0].Line,
						LineNumber: scanner.getLineNumber(offenders[0].Offender, offenders[0].Line, patchContent, to.Path(), lineLookup),
						Offender:   offenders[0].Offender,
						Rule:       rule,
						Tags:       rule.Tags,
						SecretHash: models.GenerateHash(to.Path(), commit.Hash.String(), offenders[0].Offender),
						Severity:   rule.Severity,
					})
				}
			}
		}
	}
	return leaks, err
}

func (scanner *RegexScanner) SetWhitelist(whitelist models.Whitelist) {
	scanner.Whitelist = whitelist
}

func (scanner *RegexScanner) satisfyRules(commit *object.Commit, filePath string, content string) (leaks []models.Leak) {
	for _, rule := range scanner.Rules {
		scanner.Logger.WithFields(logging.Fields{
			"condition": "secret",
			"commit":    commit.Hash.String(),
			"file":      filePath,
			"rule":      "security",
		}).Debugf("searching for `%v`", rule.Description)
		if rule.File != "" && regexp.MustCompile(rule.File).FindString(filePath) == "" {
			continue
		}
		pattern := regexp.MustCompile(rule.Pattern)
		matches := pattern.FindAllIndex([]byte(content), -1)
		if len(matches) != 0 {
			for _, match := range matches {
				line := scanner.extractLine(match[0], match[1], content)
				offender := content[match[0]:match[1]]
				groups := pattern.FindStringSubmatch(offender)
				names := pattern.SubexpNames()
				for index, group := range groups {
					if index != 0 && names[index] == "secret" {
						offender = group
						break
					}
				}
				if len(rule.Entropies) > 0 && !scanner.validateEntropy(groups, rule) {
					scanner.Logger.Debugf("entropy not satisfied on secret %s", offender)
					continue
				}
				if scanner.checkFalsePositive(filePath, line, offender) != IsPositive {
					scanner.Logger.WithFields(logging.Fields{
						"condition": "secret",
						"commit":    commit.Hash.String(),
						"rule":      "security",
					}).Warningf("false positive secret %s", offender)
					continue
				}
				leaks = append(leaks, models.Leak{
					Author:     commit.Author.Name,
					CommitHash: commit.Hash.String(),
					CreatedAt:  commit.Author.When,
					File:       filePath,
					Line:       line,
					Offender:   offender,
					Rule:       rule,
					Tags:       rule.Tags,
					SecretHash: models.GenerateHash(filePath, line, offender),
					Severity:   rule.Severity,
				})
			}
		}
	}
	return leaks
}

func (scanner *RegexScanner) checkFalsePositive(filePath string, line string, secret string) int {
	for _, isPositiveFunc := range falsePostiveHeuristics {
		status := isPositiveFunc(filePath, line, secret)
		if status != IsPositive {
			return status
		}
	}
	return IsPositive
}

func (scanner *RegexScanner) extractLine(start int, end int, content string) string {
	for start != 0 && content[start] != '\n' {
		start--
	}
	if content[start] == '\n' {
		start++
	}
	for end < len(content)-1 && content[end] != '\n' {
		end++
	}
	return content[start:end]
}

func (scanner *RegexScanner) getLineNumber(offender string, line string, patchContent string, filePath string, lineLookup map[string]bool) (lineNumber int) {
	i := strings.Index(patchContent, fmt.Sprintf("\n+++ b/%s", filePath))
	filePatchContent := patchContent[i+1:]
	i = strings.Index(filePatchContent, "diff --git")
	if i != -1 {
		filePatchContent = filePatchContent[:i]
	}
	chunkStartLine := 0
	currLine := 0
	for _, patchLine := range strings.Split(filePatchContent, "\n") {
		if strings.HasPrefix(patchLine, "@@") {
			i := strings.Index(patchLine, "+")
			pairs := strings.Split(strings.Split(patchLine[i+1:], " @@")[0], ",")
			chunkStartLine, _ = strconv.Atoi(pairs[0])
			currLine = -1
		}
		if strings.HasPrefix(patchLine, "-") {
			currLine--
		}
		if strings.HasPrefix(patchLine, "+") && strings.Contains(patchLine, line) {
			lineNumber := chunkStartLine + currLine
			if _, ok := lineLookup[fmt.Sprintf("%s%s%d%s", offender, line, lineNumber, filePath)]; !ok {
				lineLookup[fmt.Sprintf("%s%s%d%s", offender, line, lineNumber, filePath)] = true
				return lineNumber
			}
		}
		currLine++
	}
	return 1
}

func (scanner *RegexScanner) validateEntropy(groups []string, rule models.Rule) bool {
	for _, condition := range rule.Entropies {
		if len(groups) > condition.Group {
			entropy := ShannonEntropy(groups[condition.Group])
			if entropy >= condition.MinThreshold && entropy <= condition.MaxThreshold {
				return true
			}
		}
	}
	return false
}

// NewRegexScanner create new regular expression
func NewRegexScanner(logger logging.Interface, rules []models.Rule, whitelist models.Whitelist) *RegexScanner {
	return &RegexScanner{
		Logger:    logger,
		Rules:     rules,
		Whitelist: whitelist,
	}
}
