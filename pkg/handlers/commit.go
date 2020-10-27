package handlers

import (
	"context"
	"regexp"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/util"
)

// CommitHandler handle commit messages
type CommitHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (commitHandler *CommitHandler) GetType() string {
	return core.HandlerTypeCommits
}

// Handle checking commit message with defined rules
func (commitHandler *CommitHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeCommit {
		for _, condition := range rule.Conditons {
			if canSkip := core.CanSkip(commit, rule.Type, condition.Type); canSkip {
				continue
			}
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			commitHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("Processing commit rule")
			switch condition.Type {
			case "pattern":
				commitHandler.Logger.Debugf("Commit pattern `%v`", condition.Condition)
				matches := regexp.MustCompile(condition.Condition).FindAllString(commit.Message, -1)
				if len(matches) == 0 {
					// Check if we can skip this rule
					if !commitHandler.canSkip(ctx, commit.Message, condition) {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Message `{{- .Commit.Message -}}` does't satisfy condition"))
					}
				}
			case "length":
				// TODO: dynamically check operation
				predicates := make(map[string]string)
				predicates["eq"] = "!="
				predicates["ge"] = "<"
				predicates["gt"] = "<="
				predicates["le"] = ">"
				predicates["lt"] = ">="
				predicates["ne"] = "=="
				// Test message length based on "eq", "ne", "lt", "le", "ge", "gt" predicates
				messageLength := len(commit.Message)
				matches := regexp.MustCompile(string(`(eq|ge|gt|le|lt|ne)\s+([0-9]+)`)).FindStringSubmatch(condition.Condition)
				if len(matches) < 3 {
					commitHandler.Logger.Errorf("Invalid length condition `%v`", condition.Condition)
					continue
				}
				conditionLength, err := strconv.Atoi(matches[2])
				if err != nil {
					commitHandler.Logger.Errorf("Failed to parse int %v", err)
				}
				data.Operator = matches[1]
				data.Operand = matches[2]
				commitHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash.String(),
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           rule.Type,
					"user_id":        util.GetUserID(ctx),
				}).Debugf("Check if commit length %v %v %v", messageLength, matches[1], conditionLength)
				switch matches[1] {
				case "eq":
					if messageLength != conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message equal to {{ .Operand }}"))
					}
				case "ge":
					if messageLength < conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message shorter or equal than {{ .Operand }}"))
					}
				case "gt":
					if messageLength <= conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message shorter than {{ .Operand }}"))
					}
				case "le":
					if messageLength > conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message longer or equal than {{ .Operand }}"))
					}
				case "lt":
					if messageLength >= conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message longer than {{ .Operand }}"))
					}
				case "ne":
					if messageLength == conditionLength {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "Commit message not equal to {{ .Operand }}"))
					}
				default:
					commitHandler.Logger.WithFields(logging.Fields{
						"commit":         commit.Hash.String(),
						"condition":      condition.Type,
						"correlation_id": util.GetRequestID(ctx),
						"rule":           rule.Type,
						"user_id":        util.GetUserID(ctx),
					}).Infof("Unknown operation %v for length condition", matches[1])
				}
			case "email":
				matches := regexp.MustCompile(condition.Condition).FindStringSubmatch(commit.Author.Email)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "User email `{{ .Commit.Author.Email }}` does't satisfy condition"))
				}
			default:
				commitHandler.Logger.WithFields(logging.Fields{
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

func (commitHandler *CommitHandler) canSkip(ctx context.Context, commitSubject string, condition hook.Condition) bool {
	if condition.Skip != "" {
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(commitSubject)
		if len(matches) > 0 {
			commitHandler.Logger.WithFields(logging.Fields{
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Infof("Rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
