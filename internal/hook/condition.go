package hook

// ConditionType condition type
type ConditionType string

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
