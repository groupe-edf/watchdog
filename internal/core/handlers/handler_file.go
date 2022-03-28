package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime/trace"

	"github.com/dustin/go-humanize"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

var _ Handler = (*FileHandler)(nil)

// FileHandler handle committed files
type FileHandler struct {
	AbstractHandler
}

func (*FileHandler) Name() string {
	return "file"
}

// GetType return handler type
func (fileHandler *FileHandler) GetType() HandlerType {
	return HandlerTypeCommits
}

// Handle checking files with defined rules
func (fileHandler *FileHandler) Handle(ctx context.Context, commit *models.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	defer trace.StartRegion(ctx, "Scanner.Scan").End()
	trace.Log(ctx, "file", commit.Hash)
	var parentCommit string
	if len(commit.Parents) > 0 {
		parentCommit = commit.Parents[0]
	}
	entries, err := git.GetAffectedFiles(commit.Repository, &git.AffectedFilesOptions{
		DiffFilter:  "AM",
		NewCommitID: commit.Hash,
		OldCommitID: parentCommit,
	})
	if err != nil {
		return nil, err
	}
	if policy.Type == models.PolicyTypeFile {
		for _, condition := range policy.Conditions {
			if canSkip := CanSkip(commit, policy.Type, condition.Type); canSkip {
				continue
			}
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
			fileHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash,
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing file analysis")
			switch condition.Type {
			case models.ConditionTypeExtension:
				for _, entry := range entries {
					matches := regexp.MustCompile(fmt.Sprintf(`(.*?)(%s)$`, condition.Pattern)).FindAllString(entry.SrcPath, -1)
					if len(matches) != 0 {
						data.Object = entry.SrcPath
						data.Operand = filepath.Ext(entry.SrcPath)
						if !fileHandler.canSkip(ctx, entry.SrcPath, condition) {
							issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "{{ .Object }} : *.{{ .Condition.Condition }} files are not allowed"))
						}
					}
				}
			case models.ConditionTypeSize:
				matches := regexp.MustCompile(string(`(?i)(lt)\s*([0-9\s*mb|kb]+)`)).FindStringSubmatch(condition.Pattern)
				if len(matches) < 3 {
					fileHandler.Logger.Errorf("invalid file size condition %v", condition.Pattern)
					continue
				}
				threshold, err := humanize.ParseBytes(matches[2])
				if err != nil {
					fileHandler.Logger.Errorf("error parsing file size %v : %v", matches[2], err)
				}
				repository, _ := git.NewRepository(commit.Repository.Storage)
				tree := git.NewTree(repository, commit.Hash)
				for _, entry := range entries {
					blob, err := tree.GetBlobByPath(ctx, entry.SrcPath)
					if err != nil {
						fileHandler.Logger.Errorf("error loading commit tree %v : %v", tree.ID, err)
					}
					var fileSize = uint64(blob.Size(ctx))
					data.Object = entry.SrcPath
					data.Operator = matches[1]
					data.Operand = humanize.Bytes(threshold)
					data.Value = humanize.Bytes(fileSize)
					switch matches[1] {
					case "lt":
						if fileSize >= threshold {
							issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}"))
						}
					default:
						fileHandler.Logger.WithFields(logging.Fields{
							"commit":         commit.Hash,
							"condition":      condition.Type,
							"correlation_id": util.GetRequestID(ctx),
							"rule":           policy.Type,
							"user_id":        util.GetUserID(ctx),
						}).Warningf("unknown operation %v for file size condition", matches[1])
					}
				}
			default:
				fileHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash,
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           policy.Type,
					"user_id":        util.GetUserID(ctx),
				}).Warning("unsuported condition: %s", condition.Type)
			}
		}
	}
	return issues, nil
}

// CheckExtension check file extension
func (fileHandler *FileHandler) CheckExtension(fileName string, fileExtension string) (result bool, err error) {
	return true, nil
}

func (fileHandler *FileHandler) canSkip(ctx context.Context, fileName string, condition models.Condition) bool {
	if condition.Skip != "" {
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(fileName)
		if len(matches) > 0 {
			fileHandler.Logger.WithFields(logging.Fields{
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"user_id":        util.GetUserID(ctx),
			}).Infof("rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
