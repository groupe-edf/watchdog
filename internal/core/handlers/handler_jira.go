package handlers

import (
	"context"
	"os"
	"regexp"

	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/jira"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/util"
)

const (
	issueReferencePattern = "([A-Za-z-]+-[\\d]+)"
)

var (
	_       Handler = (*JiraHandler)(nil)
	baseURL         = os.Getenv("JIRA_URL")
)

// JiraHandler handle jira issues
type JiraHandler struct {
	AbstractHandler
}

// GetType return handler type
func (jiraHandler *JiraHandler) GetType() HandlerType {
	return HandlerTypeCommits
}

// Handle checking files with defined rules
func (jiraHandler *JiraHandler) Handle(ctx context.Context, commit *models.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	if policy.Type == models.PolicyTypeJira {
		for _, condition := range policy.Conditions {
			data := issue.Data{
				Commit: models.Commit{
					Author: &models.Signature{
						Email: commit.Author.Email,
						Name:  commit.Author.Name,
					},
					Hash: commit.Hash,
				},
				Condition: condition,
			}
			jiraHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing jira analysis")
			switch condition.Type {
			case models.ConditionTypeIssue:
				// Check if commit message contains issue reference
				matches := regexp.MustCompile(issueReferencePattern).FindStringSubmatch(commit.Subject)
				if len(matches) == 0 {
					var severity models.Score = models.SeverityHigh
					if jiraHandler.canSkip(commit.Subject, condition) {
						severity = models.SeverityLow
					}
					issues = append(issues, issue.NewIssue(policy, condition.Type, data, severity, "Commit message require JIRA Issue key"))
				} else {
					// Check if the issue exist by calling Jira restful API "https://jira.atlassian.com/rest/api/latest/issue/JRA-9"
					issueID := matches[1]
					if baseURL != "" {
						jiraHandler.Logger.Debugf("Issue reference verification on `%v` Jira instance", baseURL)
						jiraClient, err := jira.New(baseURL)
						if err != nil {
							jiraHandler.Logger.Debugf("Error when creating http client %v", err)
						}
						jiraClient.Authentication.SetBasicAuth(os.Getenv("JIRA_USERNAME"), os.Getenv("JIRA_PASSWORD"))
						_, err = jiraClient.GetIssue(issueID)
						if err != nil {
							issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Jira issue not found"))
						}
					}
				}
			default:
				jiraHandler.Logger.WithFields(logging.Fields{
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

func (jiraHandler *JiraHandler) canSkip(commitSubject string, condition models.Condition) bool {
	if condition.Skip != "" {
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(commitSubject)
		if len(matches) > 0 {
			jiraHandler.Logger.Infof("rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
