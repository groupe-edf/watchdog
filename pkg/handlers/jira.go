package handlers

import (
	"context"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/jira"
	"github.com/sirupsen/logrus"
)

const (
	issueReferencePattern = "([A-Za-z-]+-[\\d]+)"
)

var (
	baseURL = os.Getenv("JIRA_URL")
)

// JiraHandler handle jira issues
type JiraHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (jiraHandler *JiraHandler) GetType() string {
	return core.HandlerTypeCommits
}

// Handle checking files with defined rules
func (jiraHandler *JiraHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeJira {
		for _, condition := range rule.Conditons {
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			jiraHandler.Logger.WithFields(logrus.Fields{
				"commit":         commit.Hash,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Info("Processing jira analysis")
			switch condition.Type {
			case "issue":
				// Check if commit message contains issue reference
				matches := regexp.MustCompile(issueReferencePattern).FindStringSubmatch(commit.Message)
				if len(matches) == 0 {
					var severity issue.Score = issue.SeverityHigh
					if jiraHandler.canSkip(commit.Message, condition) {
						severity = issue.SeverityLow
					}
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, severity, "Commit message require JIRA Issue key"))
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
							issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Jira issue not found"))
						}
					}
				}
			default:
				jiraHandler.Logger.WithFields(logrus.Fields{
					"commit":         commit.Hash,
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

func (jiraHandler *JiraHandler) canSkip(commitSubject string, condition hook.Condition) bool {
	if condition.Skip != "" {
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(commitSubject)
		if len(matches) > 0 {
			jiraHandler.Logger.Infof("Rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
