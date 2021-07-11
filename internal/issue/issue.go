package issue

import (
	"bytes"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/groupe-edf/watchdog/internal/models"
)

var (
	// FunctionsMap helper functions
	FunctionsMap = template.FuncMap{
		"hide": HideSecret,
	}
)

// Data data to be used for message rendering
type Data struct {
	Branch    string
	Commit    models.Commit
	Condition models.Condition
	Object    string
	Operator  string
	Operand   string
	Tag       string
	Value     string
}

// NewIssue create new issue
func NewIssue(policy models.Policy, conditionType models.ConditionType, data Data, severity models.Score, messageTemplate string) models.Issue {
	if data.Condition.Ignore {
		severity = models.SeverityLow
	}
	if data.Condition.RejectionMessage != "" {
		messageTemplate = data.Condition.RejectionMessage
	}
	var message bytes.Buffer
	t := template.Must(template.New("").Funcs(FunctionsMap).Parse(messageTemplate))
	_ = t.Execute(&message, data)
	issue := models.Issue{
		Commit:        data.Commit,
		ConditionType: conditionType,
		Policy:        policy,
		PolicyType:    policy.Type,
		Message:       message.String(),
		Offender: &models.Offender{
			Object:   data.Object,
			Operand:  data.Operand,
			Operator: data.Operator,
			Value:    data.Value,
		},
		Severity: severity,
	}
	return issue
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
