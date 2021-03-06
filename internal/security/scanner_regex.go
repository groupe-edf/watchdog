package security

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/logging"
)

// IsFalsePositiveFunc function to check false positive secrets
type IsFalsePositiveFunc func(file string, line string, secret string) int

const (
	// IsPositive string is secret
	IsPositive int = iota
	// IsFile string is a path
	IsFile
	// IsFunction string is function
	IsFunction
	// IsPlaceholder string is placeholder
	IsPlaceholder
	// IsVariable string is variable
	IsVariable
	// PerCharThreshold entropy per character threshold
	PerCharThreshold = 3
)

var (
	falsePostiveHeuristics = []IsFalsePositiveFunc{
		// Check if is dynamic variable
		IsFalsePositive,
	}
	// SupportedLanguages list of supported languages
	SupportedLanguages = []string{"go", "groovy", "java", "js", "py"}
)

// RegexScanner data struct
type RegexScanner struct {
	Logger  logging.Interface
	Options Options
	Rules   []Rule
}

// AddAllowedFiles add files to allowed list
func (scanner *RegexScanner) AddAllowedFiles(files *regexp.Regexp) {
	scanner.Options.AllowList.Files = append(scanner.Options.AllowList.Files, files)
}

// Scan scan commit
func (scanner *RegexScanner) Scan(commit *object.Commit) (leaks []Leak, err error) {
	fileIter, err := commit.Files()
	if err != nil {
		return leaks, err
	}
	err = fileIter.ForEach(func(file *object.File) error {
		isBinary, err := file.IsBinary()
		if isBinary {
			return nil
		} else if err != nil {
			return err
		}
		// Check global allow list
		if len(scanner.Options.AllowList.Files) != 0 {
			for _, fileRegex := range scanner.Options.AllowList.Files {
				if fileRegex.FindString(file.Name) != "" {
					return nil
				}
			}
		}
		fileContent, err := file.Contents()
		if err != nil {
			return err
		}
		leaks = append(leaks, scanner.SatisfyRules(commit, file.Name, fileContent)...)
		return nil
	})
	return leaks, err
}

// SatisfyRules check all security rules
func (scanner *RegexScanner) SatisfyRules(commit *object.Commit, filePath string, content string) (leaks []Leak) {
	for _, rule := range scanner.Rules {
		scanner.Logger.WithFields(logging.Fields{
			"condition": "secret",
			"commit":    commit.Hash.String(),
			"file":      filePath,
			"rule":      "security",
		}).Debugf("searching for `%v`", rule.Description)
		if rule.File != nil && rule.File.FindString(filePath) == "" {
			continue
		}
		matches := rule.Regexp.FindAllIndex([]byte(content), -1)
		if len(matches) != 0 {
			for _, match := range matches {
				line := scanner.extractLine(match[0], match[1], content)
				offender := content[match[0]:match[1]]
				groups := rule.Regexp.FindStringSubmatch(offender)
				names := rule.Regexp.SubexpNames()
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
				file, _ := commit.File(filePath)
				reader, _ := file.Reader()
				leaks = append(leaks, Leak{
					File:       filePath,
					Line:       line,
					LineNumber: scanner.getLineNumber(line, reader),
					Offender:   offender,
					Rule:       rule.Description,
					Tags:       rule.Tags,
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

func (scanner *RegexScanner) getLineNumber(line string, reader io.Reader) (lineNumber int) {
	bufferScanner := bufio.NewScanner(reader)
	lineNumber = 1
	for bufferScanner.Scan() {
		if line == bufferScanner.Text() {
			break
		}
		lineNumber++
	}
	return lineNumber
}

func (scanner *RegexScanner) validateEntropy(groups []string, rule Rule) bool {
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
func NewRegexScanner(logger logging.Interface, options *config.Options) *RegexScanner {
	if len(options.Security.Rules) > 0 {
		var defaultRules []Rule
		for _, rule := range options.Security.Rules {
			logger.WithFields(logging.Fields{
				"condition": "secret",
				"rule":      "security",
			}).Infof("adding secret rule %s", rule.Description)
			defaultRule := NewRule(rule.Description, rule.File, rule.Regexp, rule.Severity, rule.Tags)
			if defaultRule != nil {
				defaultRules = append(defaultRules, *defaultRule)
			}
		}
		if options.Security.MergeRules {
			rules = append(rules, defaultRules...)
		} else {
			rules = defaultRules
		}
	}
	return &RegexScanner{
		Logger: logger,
		Options: Options{
			AllowList: AllowList{
				Description: "Default allowed files",
			},
		},
		Rules: rules,
	}
}

// IsFalsePositive check if secret is a false positive
func IsFalsePositive(filePath string, line string, secret string) int {
	// Secret is a variable
	if strings.HasPrefix(secret, "$") && !strings.Contains(secret[2:], "$") {
		return IsVariable
	}
	// Secret is a placeholder
	if strings.Contains(secret, "{{") || strings.Contains(secret, "}}") {
		return IsPlaceholder
	}
	// Secret is a placeholder
	if strings.HasPrefix(secret, "{") && strings.HasSuffix(secret, "}") {
		if len(secret) < 32 {
			return IsPlaceholder
		}
	}
	// Secret is a placeholder
	if strings.HasPrefix(secret, "${") && strings.HasSuffix(secret, "}") {
		return IsPlaceholder
	}
	extension := filepath.Ext(filePath)
	openBrackets := strings.Count(secret, "(")
	closeBrackets := strings.Count(secret, ")")
	// Secret is method or function
	if IsSupportedLanguage(extension) {
		if openBrackets >= 1 {
			if openBrackets == closeBrackets {
				return IsFunction
			}
		}
		if strings.Count(line, "\""+secret+"\"") < 1 {
			return IsVariable
		}
	}
	if strings.HasSuffix(secret, ";") {
		return IsVariable
	}
	return IsPositive
}

// IsSupportedLanguage check if extension is suported
func IsSupportedLanguage(language string) bool {
	for _, supported := range SupportedLanguages {
		if language == "."+supported {
			return true
		}
	}
	return false
}
