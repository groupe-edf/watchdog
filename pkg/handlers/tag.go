package handlers

import (
	"context"

	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/sirupsen/logrus"
)

// TagHandler handle tags
type TagHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (tagHandler *TagHandler) GetType() string {
	return core.HandlerTypeRefs
}

// Handle checking tags with defined rules
func (tagHandler *TagHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	// Handler must run only on tag changes
	// TODO: check only heads refs
	if rule.Type == hook.TypeTag && tagHandler.Info.RefType == "tags" {
		for _, condition := range rule.Conditons {
			data := issue.Data{
				Condition: condition,
				Tag:       tagHandler.Info.RefName,
			}
			tagHandler.Logger.WithFields(logrus.Fields{
				"tag":            data.Tag,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Info("Processing tag rule")
			switch condition.Type {
			case "semver":
				_, err := semver.NewVersion(data.Tag)
				if err != nil {
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Tag version `{{ .Tag }}` must respect semantic versionning v2.0.0 https://semver.org/"))
				}
			default:
				tagHandler.Logger.WithFields(logrus.Fields{
					"tag":            data.Tag,
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           rule.Type,
					"user_id":        util.GetUserID(ctx),
				}).Info("Unsuported condition")
			}
		}
	}
	return issues, nil
}
