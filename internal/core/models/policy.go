package models

import (
	"time"
)

// PolicyType policy type
type PolicyType string

const (
	// PolicyTypeBranch constant for branch handler
	PolicyTypeBranch PolicyType = "branch"
	// PolicyTypeCommit constant for commit handler
	PolicyTypeCommit PolicyType = "commit"
	// PolicyTypeFile constant for file handler
	PolicyTypeFile PolicyType = "file"
	// PolicyTypeJira constant for jira handler
	PolicyTypeJira PolicyType = "jira"
	// PolicyTypeSecurity constant for security handler
	PolicyTypeSecurity PolicyType = "security"
	// PolicyTypeTag constant for tag handler
	PolicyTypeTag PolicyType = "tag"
)

// Policy represents a policy
// swagger:model
type Policy struct {
	ID         int64       `json:"id,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
	// swagger:strfmt date-time
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Description string     `json:"description,omitempty"`
	DisplayName string     `json:"display_name,omitempty"`
	Enabled     bool       `json:"enabled"`
	EventType   []string   `json:"event_type,omitempty"`
	Name        string     `json:"name,omitempty"`
	Severity    Severity   `json:"severity,omitempty"`
	// Type of the policy
	//
	// type: string
	// enum: branch,commit,file,jira,security,tag
	Type PolicyType `json:"type,omitempty"`
}
