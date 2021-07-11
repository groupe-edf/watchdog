package handlers

import (
	"context"
	"regexp"
	"runtime/trace"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/util"
)

var (
	_ Handler = (*CommitHandler)(nil)
)

// CommitHandler handle commit messages
type CommitHandler struct {
	AbstractHandler
}

// GetType return handler type
func (commitHandler *CommitHandler) GetType() HandlerType {
	return HandlerTypeCommits
}

// Handle checking commit message with defined rules
func (commitHandler *CommitHandler) Handle(ctx context.Context, commit *object.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	defer trace.StartRegion(ctx, "Scanner.Scan").End()
	trace.Log(ctx, "commit", commit.Hash.String())
	if policy.Type == models.PolicyTypeCommit {
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
			commitHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing commit rule")
			switch condition.Type {
			case models.ConditionTypePattern:
				commitHandler.Logger.Debugf("Commit pattern `%v`", condition.Pattern)
				matches := regexp.MustCompile(condition.Pattern).FindAllString(commit.Message, -1)
				if len(matches) == 0 {
					// Check if we can skip this rule
					if !commitHandler.canSkip(ctx, commit.Message, condition) {
						data.Value = commit.Message
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Message `{{- .Commit.Message -}}` does not satisfy condition"))
					}
				}
			case models.ConditionTypeLength:
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
				matches := regexp.MustCompile(string(`(eq|ge|gt|le|lt|ne)\s+([0-9]+)`)).FindStringSubmatch(condition.Pattern)
				if len(matches) < 3 {
					commitHandler.Logger.Errorf("Invalid length condition `%v`", condition.Pattern)
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
					"rule":           policy.Type,
					"user_id":        util.GetUserID(ctx),
				}).Debugf("Check if commit length %v %v %v", messageLength, matches[1], conditionLength)
				switch matches[1] {
				case "eq":
					if messageLength != conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message equal to {{ .Operand }}"))
					}
				case "ge":
					if messageLength < conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message shorter or equal than {{ .Operand }}"))
					}
				case "gt":
					if messageLength <= conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message shorter than {{ .Operand }}"))
					}
				case "le":
					if messageLength > conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message longer or equal than {{ .Operand }}"))
					}
				case "lt":
					if messageLength >= conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message longer than {{ .Operand }}"))
					}
				case "ne":
					if messageLength == conditionLength {
						issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "Commit message not equal to {{ .Operand }}"))
					}
				default:
					commitHandler.Logger.WithFields(logging.Fields{
						"commit":         commit.Hash.String(),
						"condition":      condition.Type,
						"correlation_id": util.GetRequestID(ctx),
						"rule":           policy.Type,
						"user_id":        util.GetUserID(ctx),
					}).Warningf("unknown operation %v for length condition", matches[1])
				}
			case models.ConditionTypeEmail:
				matches := regexp.MustCompile(condition.Pattern).FindStringSubmatch(commit.Author.Email)
				if len(matches) == 0 {
					issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "User email `{{ .Commit.Author.Email }}` does not satisfy condition"))
				}
			default:
				commitHandler.Logger.WithFields(logging.Fields{
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

func (commitHandler *CommitHandler) canSkip(ctx context.Context, commitSubject string, condition models.Condition) bool {
	if condition.Skip != "" {
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(commitSubject)
		if len(matches) > 0 {
			commitHandler.Logger.WithFields(logging.Fields{
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Infof("rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
