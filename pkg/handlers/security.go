package handlers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/core"
	"github.com/groupe-edf/watchdog/internal/hook"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/sirupsen/logrus"
)

const (
	// BaiscAuthenticationPattern common pattern for basic authentication URl
	BaiscAuthenticationPattern string = "://[^{}[:space:]]+:(?P<secret>[^{}[:space:]]+)@"
	// Base64Pattern Base64 pattern
	Base64Pattern string = "(?P<secret>(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4}))"
	// ConnectionString pattern for connection string like ftp, database ...
	ConnectionString string = "(?P<scheme>[a-z0-9+-.]{3,30}://)?[0-9a-z-]{3,30}:" + PasswordPattern + "@(?P<host>[0-9a-z-.]{1,50})(?::(?P<port>[0-9]{1,5}))?]?"
	// PasswordExcludePattern pattern to be excluded in password pattern
	PasswordExcludePattern string = ""
	// PasswordPattern Password pattern
	// FIXME: exclude variables, no support to negative lookahead and lookbehind in golang re2
	PasswordPattern string = "(?P<secret>[a-z0-9!?$)()=<>\\/%@#*&{}_^-]{6,45})"
	// PasswordPrefixPattern token used to recognize passwords
	PasswordPrefixPattern string = "(?:(?:pass(?:w(?:or)?d)?)|(?:p(?:s)?w(?:r)?d)|secret)"
	// SecretKeyPrefixPattern token used to recognize secrets
	SecretKeyPrefixPattern string = "(?:(?:a(?:ws|ccess|p(?:i|p(?:lication)?)))|private|se(?:nsitive|cret))"
	// ConditionIP ip condition
	ConditionIP hook.ConditionType = "ip"
	// ConditionSecret secret condition
	ConditionSecret hook.ConditionType = "secret"
	// ConditionSignature secret condition
	ConditionSignature hook.ConditionType = "signature"
)

var (
	scanners = []SecretPattern{
		{
			"AWS_ACCESS_KEY",
			regexp.MustCompile("(?P<secret>AKIA[0-9A-Z]{16})"),
			make([]string, 0),
		},
		{
			"AWS_SECRET_KEY",
			regexp.MustCompile("(?P<secret>[0-9a-zA-Z/+]{40})"),
			make([]string, 0),
		},
		{
			"BASIC_AUTHENTICATION_HEADER",
			regexp.MustCompile("(?:Authorization:[[:space:]])?Basic[[:space:]]" + Base64Pattern),
			make([]string, 0),
		},
		{
			"CONFIDENTIAL",
			regexp.MustCompile("(?i)confidential"),
			make([]string, 0),
		},
		{
			"CONNECTION_STRING",
			regexp.MustCompile("(?i)" + ConnectionString),
			make([]string, 0),
		},
		{
			"GOOGLE_ACCESS_TOKEN",
			regexp.MustCompile(string(`(?P<secret>ya29.[0-9a-zA-Z_\-]{68})`)),
			make([]string, 0),
		},
		{
			"GOOGLE_API",
			regexp.MustCompile(string(`(?P<secret>AIzaSy[0-9a-zA-Z_\\-]{33})`)),
			make([]string, 0),
		},
		{
			"NPM_AUTHENTICATION",
			regexp.MustCompile("_auth[[:space:]]*=[[:space:]]*" + Base64Pattern),
			make([]string, 0),
		},
		{
			"PASSWORD",
			regexp.MustCompile("(?i)['\"]?" + PasswordPrefixPattern + "['\"]?[[:space:]]{0,20}[=:]{1,3}[[:space:]]{0,20}[@]?['\"]?" + PasswordPattern + "['\"]?"),
			make([]string, 0),
		},
		{
			"PASSWORD_XML",
			regexp.MustCompile("(?i)<" + PasswordPrefixPattern + ">(.{5,256})</" + PasswordPrefixPattern + ">"),
			make([]string, 0),
		},
		{
			"PRIVATE_KEY",
			regexp.MustCompile(string(`(\-){5}BEGIN[[:space:]]*?(RSA|OPENSSH|DSA|EC|PGP)?[[:space:]]*?PRIVATE KEY[[:space:]]*?(BLOCK)?(\-){5}.*`)),
			make([]string, 0),
		},
		{
			"SECRET_KEY",
			regexp.MustCompile("(?i)" + SecretKeyPrefixPattern + "[[:space:]_-]?(?:key|token)[[:space:]]{0,20}[=:]{1,2}[[:space:]]{0,20}['\"]?" + PasswordPattern + "(?:[[:space:];'\",]|$)"),
			make([]string, 0),
		},
		{
			"SLACK",
			regexp.MustCompile("(?P<secret>xox.-[0-9]{12}-[0-9]{12}-[0-9a-zA-Z]{24})"),
			make([]string, 0),
		},
		{
			"TWILIO",
			regexp.MustCompile("(?P<secret>55[0-9a-fA-F]{32})"),
			make([]string, 0),
		},
		{
			"TWITTER",
			regexp.MustCompile("(?P<secret>[1-9][0-9]+-[0-9a-zA-Z]{40})"),
			make([]string, 0),
		},
	}
)

// SecretPattern data struct
type SecretPattern struct {
	Provider string
	Regexp   *regexp.Regexp
	Matches  []string
}

// SecurityHandler handle committed secrets, passwords and tokens
type SecurityHandler struct {
	core.AbstractHandler
}

// GetType return handler type
func (securityHandler *SecurityHandler) GetType() string {
	return core.HandlerTypeCommits
}

// Handle checking files for secrets
func (securityHandler *SecurityHandler) Handle(ctx context.Context, commit *object.Commit, rule *hook.Rule) (issues []issue.Issue, err error) {
	if rule.Type == hook.TypeSecurity {
		for _, condition := range rule.Conditons {
			data := issue.Data{
				Commit:    commit,
				Condition: condition,
			}
			securityHandler.Logger.WithFields(logrus.Fields{
				"commit":         commit.Hash.String(),
				"condition":      condition.Type,
				"correlation_id": util.GetRequestID(ctx),
				"rule":           rule.Type,
				"user_id":        util.GetUserID(ctx),
			}).Info("Processing security analysis")
			switch condition.Type {
			case ConditionSecret:
				fileIter, err := commit.Files()
				if err != nil {
					return issues, err
				}
				err = fileIter.ForEach(func(file *object.File) error {
					isBinary, err := file.IsBinary()
					if isBinary {
						return nil
					} else if err != nil {
						return err
					}
					fileContent, err := file.Contents()
					if err != nil {
						return err
					}
					for _, scanner := range scanners {
						securityHandler.Logger.Debugf("Searching for `%v` secret", scanner.Provider)
						matches := scanner.Regexp.FindAllString(fileContent, -1)
						if matches != nil {
							// TODO: Use 1.15 new feature SubexpIndex secret = scanner.Regexp.SubexpIndex("secret")
							scanner.Matches = append(scanner.Matches, matches...)
							// FIXME: avoid logging secrets
							if !securityHandler.canSkip(file.Name, condition) && !core.CanSkip(commit, rule.Type, condition.Type) {
								data.Value = matches[0]
								data.Object = file.Name
								issues = append(issues, issue.NewIssue(rule.Type, condition.Type, data, issue.SeverityHigh, fmt.Sprintf("You're trying to commit a %s secret: {{ .Object }}:{{ Hide .Value 4 }}", scanner.Provider)))
							}
						}
					}
					return nil
				})
				return issues, err
			// TODO: implement ip and signature hooks
			case ConditionIP:
			case ConditionSignature:
			default:
				securityHandler.Logger.WithFields(logrus.Fields{
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

func (securityHandler *SecurityHandler) canSkip(fileName string, condition hook.Condition) bool {
	if condition.Skip != "" {
		securityHandler.Logger.Debugf("Skip condition `%v` found", condition.Skip)
		matches := regexp.MustCompile(condition.Skip).FindStringSubmatch(fileName)
		if len(matches) > 0 {
			securityHandler.Logger.Debugf("Rule ignored due to skip condition `%v`", condition.Skip)
			return true
		}
	}
	return false
}
