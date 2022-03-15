package handlers

import (
	"context"
	"regexp"

	"github.com/coreos/go-semver/semver"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/util"
)

var (
	_ Handler = (*TagHandler)(nil)
)

// TagHandler handle tags
type TagHandler struct {
	AbstractHandler
}

// GetType return handler type
func (tagHandler *TagHandler) GetType() HandlerType {
	return HandlerTypeRefs
}

// Handle checking tags with defined rules
func (tagHandler *TagHandler) Handle(ctx context.Context, commit *models.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	// Handler must run only on tag changes
	// TODO: check only heads refs
	if policy.Type == models.PolicyTypeTag && tagHandler.Info.Ref.IsTag() {
		for _, condition := range policy.Conditions {
			data := issue.Data{
				Condition: condition,
				Commit: models.Commit{
					Hash: tagHandler.Info.NewRev.String(),
				},
				Tag: tagHandler.Info.Ref.Short(),
			}
			tagHandler.Logger.WithFields(logging.Fields{
				"tag":            data.Tag,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing tag rule")
			switch condition.Type {
			case models.ConditionTypeSemVer:
				_, err := semver.NewVersion(data.Tag)
				if err != nil {
					issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Tag `{{ .Tag }}` must respect semantic versionning v2.0.0 https://semver.org/"))
				}
			case models.ConditionTypePattern:
				tagHandler.Logger.Debugf("Tag pattern `%v`", condition.Pattern)
				matches := regexp.MustCompile(condition.Pattern).FindAllString(data.Tag, -1)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Tag `{{ .Tag }}` does not satisfy condition"))
				}
			default:
				tagHandler.Logger.WithFields(logging.Fields{
					"tag":            data.Tag,
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
