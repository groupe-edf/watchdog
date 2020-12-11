package issue

import (
	"bytes"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/security"
)

var (
	// FunctionsMap helper functions
	FunctionsMap = template.FuncMap{
		"Hide": HideSecret,
	}
)

// Data data to be used for message rendering
type Data struct {
	Branch    string
	Commit    *object.Commit
	Condition hook.Condition
	Object    string
	Operator  string
	Operand   string
	Tag       string
	Value     string
}

// Issue analysis issue
type Issue struct {
	Author    string             `json:"author"`
	Commit    string             `json:"commit"`
	Condition hook.ConditionType `json:"condition"`
	Email     string             `json:"email"`
	Handler   hook.HandlerType   `json:"handler"`
	Leaks     []security.Leak    `json:"leaks,omitempty"`
	Message   string             `json:"message"`
	Severity  Score              `json:"severity"`
}

// WithLeak attach leaks to issue
func (issue *Issue) WithLeak(leak security.Leak) {
	issue.Leaks = append(issue.Leaks, leak)
}

// WithLeaks attach leaks to issue
func (issue *Issue) WithLeaks(leaks []security.Leak) {
	issue.Leaks = leaks
}

// NewIssue create new issue
func NewIssue(handlerType hook.HandlerType, conditionType hook.ConditionType, data Data, severity Score, messageTemplate string) Issue {
	if data.Condition.Ignore {
		severity = SeverityLow
	}
	if data.Condition.RejectionMessage != "" {
		messageTemplate = data.Condition.RejectionMessage
	}
	var message bytes.Buffer
	t := template.Must(template.New("").Funcs(FunctionsMap).Parse(messageTemplate))
	_ = t.Execute(&message, data)
	return Issue{
		Author:    data.Commit.Author.Name,
		Commit:    data.Commit.Hash.String(),
		Condition: conditionType,
		Email:     data.Commit.Author.Email,
		Handler:   handlerType,
		Message:   message.String(),
		Severity:  severity,
	}
}

// HideSecret hide leaks in text
func HideSecret(secret string, reveal int) string {
	if reveal > len(secret) {
		reveal = len(secret)
	}
	if reveal == 0 {
		secret = strings.Repeat("*", utf8.RuneCountInString(secret))
	}
	if reveal > 0 {
		secret = strings.Replace(secret, secret[reveal:], strings.Repeat("*", len(secret)-reveal), 1)
	}
	return secret
}
