package handlers

import (
	"context"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/security"
	"github.com/groupe-edf/watchdog/internal/util"
)

const (
	// ConditionIP ip condition
	ConditionIP hook.ConditionType = "ip"
	// ConditionSecret secret condition
	ConditionSecret hook.ConditionType = "secret"
	// ConditionSignature secret condition
	ConditionSignature hook.ConditionType = "signature"
)

// SecurityHandler handle committed secrets, passwords and tokens
type SecurityHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (securityHandler *SecurityHandler) GetType() string {
	return core.HandlerTypeCommits
}

// Handle checking files for secrets
func (securityHandler *SecurityHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeSecurity {
		for _, condition := range rule.Conditons {
			if canSkip := core.CanSkip(commit, rule.Type, condition.Type); canSkip {
				continue
			}
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			securityHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Info("Processing security analysis")
			switch condition.Type {
			case ConditionSecret:
				// Create a new regex scanner
				options := security.Options{
					AllowList: security.AllowList{
						Files: []*regexp.Regexp{
							regexp.MustCompile("(?i)(css)$"),
						},
					},
				}
				if condition.Skip != "" {
					options.AllowList.Files = append(options.AllowList.Files, regexp.MustCompile(condition.Skip))
				}
				scanner := security.NewRegexScanner(securityHandler.Logger, options)
				leaks, err := scanner.Scan(commit)
				if err != nil {
					return nil, err
				}
				if len(leaks) > 0 {
					for _, leak := range leaks {
						data.Value = leak.Offender
						data.Object = leak.File
						issue := issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Secrets, token and passwords are forbidden, `{{ .Object }}:{{ Hide .Value 4 }}`")
						issue.WithLeak(leak)
						issues = append(issues, issue)
					}
				}
				return issues, err
			// TODO: implement ip and signature hooks
			case ConditionIP:
			case ConditionSignature:
			default:
				securityHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash.String(),
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

func (securityHandler *SecurityHandler) canSkip(fileName string, condition hook.Condition) bool {
	if condition.Skip != "" {
		securityHandler.Logger.Debugf("Skip condition `%v` found", condition.Skip)
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(fileName)
		if len(matches) > 0 {
			securityHandler.Logger.Debugf("Rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
