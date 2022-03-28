package handlers

import (
	"context"

	"github.com/gitleaks/go-gitdiff/gitdiff"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/security"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

var (
	_ Handler = (*SecurityHandler)(nil)
)

// SecurityHandler handle committed secrets, passwords and tokens
type SecurityHandler struct {
	AbstractHandler
	Scanner security.Scanner
}

func (*SecurityHandler) Name() string {
	return "security"
}

// GetType return handler type
func (securityHandler *SecurityHandler) GetType() HandlerType {
	return HandlerTypeRefs
}

// Handle checking files for secrets
func (securityHandler *SecurityHandler) Handle(ctx context.Context, commit *models.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	if policy.Type == models.PolicyTypeSecurity {
		for _, condition := range policy.Conditions {
			data := issue.Data{
				Condition: condition,
			}
			securityHandler.Logger.WithFields(logging.Fields{
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
				repository, err := git.NewRepository(securityHandler.Repository.Path())
				if err != nil {
					return nil, err
				}
				var diffFiles <-chan *gitdiff.File
				diffFiles, err = git.GetLog(ctx, []string{}, repository)
				if err != nil {
					return nil, err
				}
				for diffFile := range diffFiles {
					diffFile := diffFile
					data.Commit = models.Commit{
						Hash: diffFile.PatchHeader.SHA,
						Author: &models.Signature{
							Name:  diffFile.PatchHeader.Author.Name,
							Email: diffFile.PatchHeader.Author.Email,
						},
					}
					if diffFile.IsBinary || diffFile.IsDelete {
						continue
					}
					for _, textFragment := range diffFile.TextFragments {
						fragment := models.Fragment{
							Raw:       textFragment.Raw(gitdiff.OpAdd),
							CommitSHA: diffFile.PatchHeader.SHA,
							FilePath:  diffFile.NewName,
						}
						leaks, err := securityHandler.Scanner.Scan(fragment)
						if err != nil {
							return nil, err
						}
						if len(leaks) > 0 {
							for _, leak := range leaks {
								if diffFile.PatchHeader.Author != nil {
									leak.AuthorName = diffFile.PatchHeader.Author.Name
									leak.AuthorEmail = diffFile.PatchHeader.Author.Email
								}
								leak.CreatedAt = diffFile.PatchHeader.AuthorDate
								offender := issue.HideSecret(leak.Offender, securityHandler.Options.Security.RevealSecrets)
								securityHandler.Logger.WithFields(logging.Fields{
									"commit":         data.Commit.Hash,
									"condition":      condition.Type,
									"correlation_id": util.GetRequestID(ctx),
									"rule":           policy.Type,
									"user_id":        util.GetUserID(ctx),
								}).Debugf("potential %s secret leaked in file %s line %d: %s", leak.Rule.DisplayName, leak.File, leak.LineNumber, offender)
								data.Value = offender
								data.Object = leak.File
								issue := issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Secrets, token and passwords are forbidden, `{{ .Object }}:{{ hide .Value 4 }}`")
								issue.WithLeak(leak)
								issue.Policy = policy
								issues = append(issues, issue)
							}
						}
					}
				}
				return issues, err
			case models.ConditionTypeIP:
			case models.ConditionTypeSignature:
			default:
				securityHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash,
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
