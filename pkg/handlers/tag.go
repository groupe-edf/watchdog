package handlers

import (
	"context"
	"regexp"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/util"
)

// TagHandler handle tags
type TagHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (tagHandler *TagHandler) GetType() core.HandlerType {
	return core.HandlerTypeRefs
}

// Handle checking tags with defined rules
func (tagHandler *TagHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	// Handler must run only on tag changes
	// TODO: check only heads refs
	if rule.Type == hook.TypeTag && tagHandler.Info.Ref.IsTag() {
		for _, condition := range rule.Conditions {
			data := issue.Data{
				Condition: condition,
				Commit:    tagHandler.Info.NewRev,
				Tag:       tagHandler.Info.Ref.Short(),
			}
			tagHandler.Logger.WithFields(logging.Fields{
				"tag":            data.Tag,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing tag rule")
			switch condition.Type {
			case hook.ConditionSemVer:
				_, err := semver.NewVersion(data.Tag)
				if err != nil {
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Tag `{{ .Tag }}` must respect semantic versionning v2.0.0 https://semver.org/"))
				}
			case hook.ConditionPattern:
				tagHandler.Logger.Debugf("Tag pattern `%v`", condition.Condition)
				matches := regexp.MustCompile(condition.Condition).FindAllString(data.Tag, -1)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Tag `{{ .Tag }}` does not satisfy condition"))
				}
			default:
				tagHandler.Logger.WithFields(logging.Fields{
					"tag":            data.Tag,
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
