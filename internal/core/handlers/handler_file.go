package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/c2h5oh/datasize"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/util"
)

var (
	_ Handler = (*FileHandler)(nil)
)

// FileHandler handle committed files
type FileHandler struct {
	AbstractHandler
}

// GetType return handler type
func (fileHandler *FileHandler) GetType() HandlerType {
	return HandlerTypeCommits
}

// Handle checking files with defined rules
func (fileHandler *FileHandler) Handle(ctx context.Context, commit *object.Commit, policy models.Policy, whitelist models.Whitelist) (issues []models.Issue, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if policy.Type == models.PolicyTypeFile {
		if len(commit.ParentHashes) == 0 {
			return
		}
		parent, err := commit.Parent(0)
		if err != nil {
			return nil, err
		}
		patch, err := parent.Patch(commit)
		if err != nil {
			return nil, err
		}
		for _, condition := range policy.Conditions {
			if canSkip := CanSkip(commit, policy.Type, condition.Type); canSkip {
				continue
			}
			data := issue.Data{
				Commit: models.Commit{
					Author: commit.Author.Name,
					Email:  commit.Author.Email,
					Hash:   commit.Hash.String(),
				},
				Condition: condition,
			}
			fileHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           policy.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing file analysis")
			switch condition.Type {
			case models.ConditionTypeExtension:
				for _, filePatch := range patch.FilePatches() {
					_, to := filePatch.Files()
					if to != nil {
						fileHandler.Logger.Debugf("Checking file %v", to.Path())
						matches := regexp.MustCompile(fmt.Sprintf(`(.*?)(%s)$`, condition.Pattern)).FindAllString(to.Path(), -1)
						if len(matches) != 0 {
							data.Object = to.Path()
							data.Operand = filepath.Ext(to.Path())
							if !fileHandler.canSkip(ctx, to.Path(), condition) {
								issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "{{ .Object }} : *.{{ .Condition.Condition }} files are not allowed"))
							}
						}
					}
				}
			case models.ConditionTypeSize:
				matches := regexp.MustCompile(string(`(?i)(lt)\s*([0-9\s*mb|kb]+)`)).FindStringSubmatch(condition.Pattern)
				if len(matches) < 3 {
					fileHandler.Logger.Errorf("Invalid file size condition %v", condition.Pattern)
					continue
				}
				var size datasize.ByteSize
				err := size.UnmarshalText([]byte(matches[2]))
				if err != nil {
					fileHandler.Logger.Errorf("Error parsing file size %v : %v", matches[2], err)
				}
				for _, filePatch := range patch.FilePatches() {
					_, to := filePatch.Files()
					if to != nil {
						file, _ := commit.File(to.Path())
						var fileSize datasize.ByteSize
						err = fileSize.UnmarshalText([]byte(strconv.FormatInt(file.Size, 10)))
						if err != nil {
							fileHandler.Logger.Errorf("Error parsing file size %v : %v", matches[2], err)
						}
						data.Object = file.Name
						data.Operator = matches[1]
						data.Operand = size.HumanReadable()
						data.Value = fileSize.HumanReadable()
						switch matches[1] {
						case "lt":
							if fileSize.Bytes() >= size.Bytes() {
								issues = append(issues, issue.NewIssue(policy, condition.Type, data, models.SeverityHigh, "File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}"))
							}
						default:
							fileHandler.Logger.WithFields(logging.Fields{
								"commit":         commit.Hash.String(),
								"condition":      condition.Type,
								"correlation_id": util.GetRequestID(ctx),
								"rule":           policy.Type,
								"user_id":        util.GetUserID(ctx),
							}).Warningf("unknown operation %v for file size condition", matches[1])
						}
					}
				}
			default:
				fileHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash.String(),
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
