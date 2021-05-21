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
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/util"
)

// FileHandler handle committed files
type FileHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (fileHandler *FileHandler) GetType() core.HandlerType {
	return core.HandlerTypeCommits
}

// Handle checking files with defined rules
func (fileHandler *FileHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeFile {
		fileIter, err := commit.Files()
		if err != nil {
			fileHandler.Logger.Fatalf("error when loading commit files, %v", err)
		}
		var files []*object.File
		err = fileIter.ForEach(func(file *object.File) error {
			files = append(files, file)
			return nil
		})
		if err != nil {
			fileHandler.Logger.Fatalf("error when loading commit files, %v", err)
		}
		for _, condition := range rule.Conditions {
			if canSkip := core.CanSkip(commit, rule.Type, condition.Type); canSkip {
				continue
			}
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			fileHandler.Logger.WithFields(logging.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Debug("processing file analysis")
			switch condition.Type {
			case hook.ConditionExtension:
				for _, file := range files {
					fileHandler.Logger.Debugf("Checking file %v", file.Name)
					matches := regexp.MustCompile("(.+)."+condition.Condition).FindAllString(file.Name, -1)
					if len(matches) != 0 {
						data.Object = file.Name
						if !fileHandler.canSkip(ctx, file.Name, condition) {
							issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, "{{ .Object }} : *.{{ .Condition.Condition }} files are not allowed"))
						}
					}
				}
			case hook.ConditionSize:
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
						fileHandler.Logger.WithFields(logging.Fields{
							"commit":         commit.Hash.String(),
							"condition":      condition.Type,
							"correlation_id": util.GetRequestID(ctx),
							"rule":           rule.Type,
							"user_id":        util.GetUserID(ctx),
						}).Warningf("unknown operation %v for file size condition", matches[1])
					}
				}
			default:
				fileHandler.Logger.WithFields(logging.Fields{
					"commit":         commit.Hash.String(),
					"condition":      condition.Type,
					"correlation_id": util.GetRequestID(ctx),
					"rule":           rule.Type,
					"user_id":        util.GetUserID(ctx),
				}).Warning("unsuported condition")
			}
		}
	}
	return issues, nil
}

// CheckExtension check file extension
func (fileHandler *FileHandler) CheckExtension(fileName string, fileExtension string) (result bool, err error) {
	return true, nil
}

func (fileHandler *FileHandler) canSkip(ctx context.Context, fileName string, condition hook.Condition) bool {
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
