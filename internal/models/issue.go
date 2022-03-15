package models

// Issue analysis issue
type Issue struct {
	ID            string        `json:"id,omitempty"`
	Commit        Commit        `json:"commit"`
	ConditionType ConditionType `json:"condition_type"`
	Leaks         []Leak        `json:"leaks,omitempty"`
	Message       string        `json:"message"`
	Offender      *Offender     `json:"offender,omitempty"`
	Policy        Policy        `json:"policy"`
	PolicyType    PolicyType    `json:"policy_type"`
	Repository    *Repository   `json:"repository,omitempty"`
	Severity      Score         `json:"severity"`
}

type Offender struct {
	Object   string `json:"object,omitempty"`
	Operator string `json:"operator,omitempty"`
	Operand  string `json:"operand,omitempty"`
	Value    string `json:"value,omitempty"`
}

// WithCommit
func (issue *Issue) WithCommit(commit Commit) {
	issue.Commit = commit
}

// WithLeak attach leaks to issue
func (issue *Issue) WithLeak(leak Leak) {
	issue.Leaks = append(issue.Leaks, leak)
}

// WithLeaks attach leaks to issue
func (issue *Issue) WithLeaks(leaks []Leak) {
	issue.Leaks = leaks
}

func NewIssue() *Issue {
	return &Issue{}
}
