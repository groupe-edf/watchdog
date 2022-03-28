package models

// ConditionType condition type
type ConditionType string

const (
	// ConditionTypeEmail email condition
	ConditionTypeEmail ConditionType = "email"
	// ConditionTypeExtension file extension condition
	ConditionTypeExtension ConditionType = "extension"
	// ConditionTypeIP ip condition
	ConditionTypeIP ConditionType = "ip"
	// ConditionTypeIssue issue condition
	ConditionTypeIssue ConditionType = "issue"
	// ConditionTypeLength length condition
	ConditionTypeLength ConditionType = "length"
	// ConditionTypePattern pattern condition
	ConditionTypePattern ConditionType = "pattern"
	// ConditionTypeProtected protected condition
	ConditionTypeProtected ConditionType = "protected"
	// ConditionTypeSecret secret condition
	ConditionTypeSecret ConditionType = "secret"
	// ConditionTypeSemVer semver condition
	ConditionTypeSemVer ConditionType = "semver"
	// ConditionTypeSignature secret condition
	ConditionTypeSignature ConditionType = "signature"
	// ConditionTypeSize file size condition
	ConditionTypeSize ConditionType = "size"
)

// ConditionType policy condition struct
type Condition struct {
	ID               int64         `json:"id,omitempty"`
	Ignore           bool          `json:"ignore,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	PolicyID         int64         `json:"policy_id,omitempty"`
	RejectionMessage string        `json:"rejection_message,omitempty"`
	Skip             string        `json:"skip,omitempty"`
	Type             ConditionType `json:"type,omitempty"`
}
