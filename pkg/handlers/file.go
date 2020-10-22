package handlers

import (
	"context"
	"regexp"
	"strconv"

	"github.com/c2h5oh/datasize"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/sirupsen/logrus"
)

const (
	// ConditionExtension file extension condition
	ConditionExtension hook.ConditionType = "extension"
	// ConditionSize file size condition
	ConditionSize hook.ConditionType = "size"
)

// FileHandler handle committed files
type FileHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (fileHandler *FileHandler) GetType() string {
	return core.HandlerTypeCommits
}

// Handle checking files with defined rules
func (fileHandler *FileHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeFile {
		fileIter, _ := commit.Files()
		var files []*object.File
		_ = fileIter.ForEach(func(file *object.File) error {
			files = append(files, file)
			return nil
		})
		if err != nil {
			fileHandler.Logger.Fatalf("GetFiles() error when loading commit files, %v", err)
		}
		for _, condition := range rule.Conditons {
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			fileHandler.Logger.WithFields(logrus.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Info("Processing file analysis")
			switch condition.Type {
			case ConditionExtension:
				for _, file := range files {
					fileHandler.Logger.Debugf("Checking file %v", file.Name)
					matches := regexp.MustCompile("(.+)."+condition.Condition).FindAllString(file.Name, -1)
					if len(matches) != 0 {
						issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "*.{{ .Condition.Condition }} files are not allowed"))
					}
				}
			case ConditionSize:
				matches := regexp.MustCompile(string(`(?i)(lt)\s*([0-9\s*mb|kb]+)`)).FindStringSubmatch(condition.Condition)
				if len(matches) < 3 {
					fileHandler.Logger.Errorf("Invalid file size condition %v", condition.Condition)
					continue
				}
				var size datasize.ByteSize
				err := size.UnmarshalText([]byte(matches[2]))
				if err != nil {
					fileHandler.Logger.Errorf("Error parsing file size %v : %v", matches[2], err)
				}
				for _, file := range files {
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
							issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}"))
						}
					default:
						fileHandler.Logger.WithFields(logrus.Fields{
							"commit":         commit.Hash.String(),
							"condition":      condition.Type,
							"correlation_id": util.GetRequestID(ctx),
							"rule":           rule.Type,
							"user_id":        util.GetUserID(ctx),
						}).Infof("Unknown operation %v for file size condition", matches[1])
					}
				}
			default:
				fileHandler.Logger.WithFields(logrus.Fields{
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

// CheckExtension check file extension
func (fileHandler *FileHandler) CheckExtension(fileName string, fileExtension string) (result bool, err error) {
	return true, nil
}
