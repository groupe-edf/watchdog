package handlers

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/security"
	"github.com/groupe-edf/watchdog/internal/util"
)

var (
	_ Handler = (*SecurityHandler)(nil)
)

// SecurityHandler handle committed secrets, passwords and tokens
type SecurityHandler struct {
	AbstractHandler
	Scanner security.Scanner
}

// GetType return handler type
func (securityHandler *SecurityHandler) GetType() HandlerType {
	return HandlerTypeCommits
}

// Handle checking files for secrets
func (securityHandler *SecurityHandler) Handle(ctx context.Context, commit *object.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	if policy.Type == models.PolicyTypeSecurity {
		for _, condition := range policy.Conditions {
			if canSkip := CanSkip(commit, policy.Type, condition.Type); canSkip {
				continue
			}
			data := issue.Data{
				Commit: models.Commit{
					Author:  commit.Author.Name,
					Email:   commit.Author.Email,
					Hash:    commit.Hash.String(),
					Message: strings.TrimSuffix(commit.Message, "\n"),
				},
				Condition: condition,
			}
			securityHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing security analysis")
			switch condition.Type {
			case models.ConditionTypeSecret:
				if condition.Skip != "" {
					securityHandler.Scanner.SetWhitelist(models.Whitelist{
						Files: models.List{condition.Skip},
					})
				}
				leaks, err := securityHandler.Scanner.Scan(commit)
				if err != nil {
					return nil, err
				}
				if len(leaks) > 0 {
					for _, leak := range leaks {
						offender := issue.HideSecret(leak.Offender, securityHandler.Options.Security.RevealSecrets)
						securityHandler.Logger.WithFields(logging.Fields{
							"commit":         commit.Hash.String(),
							"condition":      condition.Type,
							"correlation_id": util.GetRequestID(ctx),
							"rule":           policy.Type,
							"user_id":        util.GetUserID(ctx),
						}).Debugf("potential %s secret leaked in file %s line %d: %s", leak.Rule.DisplayName, leak.File, leak.LineNumber, offender)
						data.Value = offender
						data.Object = leak.File
						issue := issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Secrets, token and passwords are forbidden")
						issue.WithLeak(leak)
						issue.Policy = policy
						issues = append(issues, issue)
					}
				}
				return issues, err
			case models.ConditionTypeIP:
			case models.ConditionTypeSignature:
			default:
				securityHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash.String(),
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
