package issue

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/hook"
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
	Hash      string             `json:"hash"`
	Message   string             `json:"message"`
	Severity  Score              `json:"severity"`
	Handler   hook.HandlerType   `json:"handler"`
	Condition hook.ConditionType `json:"condition"`
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
	functionsMap := template.FuncMap{
		"Hide": func(value string, characters int) string {
			return strings.Replace(value, value[characters:], strings.Repeat("#", len(value)-characters), 1)
		},
	}
	t := template.Must(template.New("").Funcs(functionsMap).Parse(messageTemplate))
	_ = t.Execute(&message, data)
	return Issue{
		Hash:      data.Commit.Hash.String(),
		Message:   message.String(),
		Severity:  severity,
		Handler:   handlerType,
		Condition: conditionType,
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
