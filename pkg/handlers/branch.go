package handlers

import (
	"context"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/util"
)

// BranchHandler handle branch naming
type BranchHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (branchHandler *BranchHandler) GetType() string {
	return core.HandlerTypeRefs
}

// Handle chencking branch naming convention
func (branchHandler *BranchHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	// Handler must run only on branch changes
	if rule.Type == hook.TypeBranch && branchHandler.Info.RefType == "heads" {
		for _, condition := range rule.Conditions {
			data := issue.Data{
				Branch:    branchHandler.Info.RefName,
				Commit:    branchHandler.Info.NewRev,
				Condition: condition,
			}
			branchHandler.Logger.WithFields(logging.Fields{
				"branch":         data.Branch,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing branch rule")
			switch condition.Type {
			case "pattern":
				// User created new branch, check naming convention
				matches := regexp.MustCompile(condition.Condition).FindAllString(data.Branch, -1)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Branch name `{{ .Branch }}` does not satisfy condition"))
				}
			case "protected":
				// Reject push if the user want to delete a protected branch
				if branchHandler.Info.RefType == "heads" {
					matches := regexp.MustCompile(condition.Condition).FindStringSubmatch(data.Branch)
					if len(matches) > 0 {
						// User try to delete protected branch
						if branchHandler.Info.NewRev.Hash == plumbing.ZeroHash {
							issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "You can't delete protected branch {{ .Branch }}"))
						}
					}
				}
			default:
				branchHandler.Logger.WithFields(logging.Fields{
					"branch":         data.Branch,
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           rule.Type,
					"user_id":        util.GetUserID(ctx),
				}).Warning("unsuported condition")
			}
		}
	}
	return issues, nil
}
