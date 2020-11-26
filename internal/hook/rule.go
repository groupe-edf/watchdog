package hook

// HandlerType handler type
type HandlerType string

const (
	// TypeBranch constant for branch handler
	TypeBranch HandlerType = "branch"
	// TypeCommit constant for commit handler
	TypeCommit HandlerType = "commit"
	// TypeFile constant for file handler
	TypeFile HandlerType = "file"
	// TypeJira constant for jira handler
	TypeJira HandlerType = "jira"
	// TypeSecurity constant for security handler
	TypeSecurity HandlerType = "security"
	// TypeTag constant for tag handler
	TypeTag HandlerType = "tag"
)

// Rule data structure
// Use similar mechanism as Gitlab, see https://docs.gitlab.com/ee/push_rules/push_rules.html
type Rule struct {
	Conditions  []Condition `yaml:"conditions,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enabled     bool        `yaml:"enabled,omitempty"`
	Type        HandlerType `yaml:"type,omitempty"`
}
