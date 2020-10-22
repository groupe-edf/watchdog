package handlers

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
)

// DefaultHandler handle repository
type DefaultHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (defaultHandler *DefaultHandler) GetType() string {
	return core.HandlerTypeRefs
}

// Handle chencking branch naming convention
func (defaultHandler *DefaultHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	locked := false
	if locked {
		rejectionMessage := "\n\nYou are attempting to push to the repository which has been made read-only" +
			"\nAccess denied, push blocked. Please contact the repository administrator. %s"
		data := issue.Data{
			Commit: commit,
			Condition: hook.Condition{
				RejectionMessage: rejectionMessage,
			},
		}
		issues = append(issues, issue.NewIssue(rule.Type, "", data, issue.SeverityHigh, ""))
	}
	return issues, nil
}
