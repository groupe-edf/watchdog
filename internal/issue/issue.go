package issue

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/security"
)

// Score type used by severity and confidence values
type Score int

const (
	// SeverityLow severity or confidence
	SeverityLow Score = iota
	// SeverityMedium severity or confidence
	SeverityMedium
	// SeverityHigh severity or confidence
	SeverityHigh
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

func (score Score) String() string {
	switch score {
	case SeverityHigh:
		return "high"
	case SeverityMedium:
		return "medium"
	case SeverityLow:
		return "low"
	}
	return "undefined"
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

// ParseScore parse score from string input
func ParseScore(score string) Score {
	switch score {
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	}
	return SeverityLow
}

// HideSecret hide leaks in text
func HideSecret(value string, characters int) string {
	return strings.Replace(value, value[characters:], strings.Repeat("#", len(value)-characters), 1)
}
