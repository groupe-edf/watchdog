package handlers

import (
	"context"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/util"
)

var (
	_ Handler = (*BranchHandler)(nil)
)

// BranchHandler handle branch naming
type BranchHandler struct {
	AbstractHandler
}

// GetType return handler type
func (branchHandler *BranchHandler) GetType() HandlerType {
	return HandlerTypeRefs
}

// Handle chencking branch naming convention
func (branchHandler *BranchHandler) Handle(ctx context.Context, commit *object.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	// Handler must run only on branch changes
	if policy.Type == models.PolicyTypeBranch && branchHandler.Info.Ref.IsBranch() {
		for _, condition := range policy.Conditions {
			data := issue.Data{
				Branch: branchHandler.Info.Ref.Short(),
				Commit: models.Commit{
					Hash: branchHandler.Info.NewRev.String(),
				},
				Condition: condition,
			}
			branchHandler.Logger.WithFields(logging.Fields{
				"branch":         data.Branch,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing branch rule")
			switch condition.Type {
			case models.ConditionTypePattern:
				// User created new branch, check naming convention
				matches := regexp.MustCompile(condition.Pattern).FindAllString(data.Branch, -1)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Branch name `{{ .Branch }}` does not satisfy condition"))
				}
			case models.ConditionTypeProtected:
				// Reject push if the user want to delete a protected branch
				matches := regexp.MustCompile(condition.Pattern).FindStringSubmatch(data.Branch)
				if len(matches) > 0 {
					// User try to delete protected branch
					if branchHandler.Info.NewRev.Hash == plumbing.ZeroHash {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "You can't delete protected branch {{ .Branch }}"))
					}
				}
			default:
				branchHandler.Logger.WithFields(logging.Fields{
					"branch":         data.Branch,
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           policy.Type,
					"user_id":        util.GetUserID(ctx),
				}).Warning("unsuported condition")
			}
		}
	}
	return issues, nil
}
