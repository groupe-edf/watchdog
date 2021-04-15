package hook

// ConditionType condition type
type ConditionType string

const (
	// ConditionEmail email condition
	ConditionEmail ConditionType = "email"
	// ConditionExtension file extension condition
	ConditionExtension ConditionType = "extension"
	// ConditionIP ip condition
	ConditionIP ConditionType = "ip"
	// ConditionIssue issue condition
	ConditionIssue ConditionType = "issue"
	// ConditionLength length condition
	ConditionLength ConditionType = "length"
	// ConditionPattern pattern condition
	ConditionPattern ConditionType = "pattern"
	// ConditionPattern protected condition
	ConditionProtected ConditionType = "protected"
	// ConditionSecret secret condition
	ConditionSecret ConditionType = "secret"
	// ConditionSemVer semver condition
	ConditionSemVer ConditionType = "semver"
	// ConditionSignature secret condition
	ConditionSignature ConditionType = "signature"
	// ConditionSize file size condition
	ConditionSize ConditionType = "size"
)

// Condition used in hook's rules
type Condition struct {
	Condition        string `yaml:"condition"`
	FailFast         string `yaml:"fail_fast"`
	Ignore           bool   `yaml:"ignore"`
	Only             []string
	RejectionMessage string        `yaml:"rejection_message"`
	Skip             string        `yaml:"skip"`
	Type             ConditionType `yaml:"type"`
}
