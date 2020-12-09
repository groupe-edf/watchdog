package security

import (
	"regexp"
)

// Severity rule severity
type Severity string

const (
	// BaiscAuthenticationPattern common pattern for basic authentication URl
	BaiscAuthenticationPattern string = "://[^{}[:space:]]+:(?P<secret>[^{}[:space:]]+)@"
	// Base64Pattern Base64 pattern
	Base64Pattern string = string(`(?P<secret>(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4}))`)
	// ConnectionString pattern for connection string like ftp, database ...
	ConnectionString string = "(?P<scheme>[a-z0-9+-.]{3,30}://)?[0-9a-z-]{3,30}:" + PasswordPattern + "@(?P<host>[0-9a-z-.]{1,50})(?::(?P<port>[0-9]{1,5}))?]?"
	// PasswordExcludePattern pattern to be excluded in password pattern
	PasswordExcludePattern string = ""
	// PasswordPattern Password pattern
	// FIXME: exclude variables, no support to negative lookahead and lookbehind in golang re2
	PasswordPattern string = string(`(?P<secret>[a-zA-Z0-9!?$)(.=<>\/%@#*&{}_^+-]{6,45})`)
	// PasswordPrefixPattern token used to recognize passwords
	PasswordPrefixPattern string = "(?:(?:pass(?:w(?:or)?d)?)|(?:p(?:s)?w(?:r)?d)|secret)"
	// SecretKeyPrefixPattern token used to recognize secrets
	SecretKeyPrefixPattern string = "(?:(?:a(?:ws|ccess|p(?:i|p(?:lication)?)))|private|se(?:nsitive|cret))"
	// SeverityBlocker blocker severity
	SeverityBlocker Severity = "BLOCKER"
	// SeverityCritical critical severity
	SeverityCritical Severity = "CRITICAL"
	// SeverityInfo info severity
	SeverityInfo Severity = "INFO"
	// SeverityMajor major severity
	SeverityMajor Severity = "MAJOR"
	// SeverityMinor minor severity
	SeverityMinor Severity = "MINOR"
)

var (
	rules = []Rule{
		{
			Description: "ASYMMETRIC_PRIVATE_KEY",
			Regexp:      regexp.MustCompile(string(`(\-){5}BEGIN[[:blank:]]*?(RSA|OPENSSH|DSA|EC|PGP)?[[:blank:]]*?PRIVATE[[:blank:]]KEY[[:blank:]]*?(BLOCK)?(\-){5}.*`)),
			Tags:        []string{"key"},
			Severity:    SeverityMajor,
		},
		{
			Description: "AWS_ACCESS_KEY",
			Regexp:      regexp.MustCompile("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}"),
			Tags:        []string{"aws"},
			Severity:    SeverityMajor,
		},
		{
			Description: "AWS_SECRET_KEY",
			Regexp:      regexp.MustCompile(string(`(?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]`)),
			Tags:        []string{"aws"},
			Severity:    SeverityBlocker,
		},
		{
			Description: "BASE_64",
			Regexp:      regexp.MustCompile("(?i)['\"]?((?:[_]?auth(?:Token|orization:[[:blank:]]Basic)?)['\"]?[[:blank:]=:]{1})[[:blank:]]*['\"]?" + Base64Pattern + "['\"]?"),
			Tags:        []string{"authentication", "base64"},
			Severity:    SeverityMinor,
		},
		{
			Description: "CONFIDENTIAL",
			Regexp:      regexp.MustCompile("(?i)CONFIDENTIAL"),
			Severity:    SeverityInfo,
		},
		{
			Description: "CONNECTION_STRING",
			Regexp:      regexp.MustCompile("(?i)" + ConnectionString),
			Severity:    SeverityMajor,
		},
		{
			Description: "ENTROPY",
			File:        regexp.MustCompile("(?i).*.sql$"),
			Entropies: []Entropy{
				{
					MaxThreshold: 8.0,
					MinThreshold: 4.0,
				},
			},
			Regexp:   regexp.MustCompile(string(`[0-9a-zA-Z-_!{}$.\/=]{8,120}`)),
			Tags:     []string{"entropy"},
			Severity: SeverityInfo,
		},
		{
			Description: "HTPASSWD",
			File:        regexp.MustCompile("(?i).htpasswd$"),
			Regexp:      regexp.MustCompile("(?i)[0-9a-zA-Z-_!{}$.=]{4,120}:" + PasswordPattern),
			Severity:    SeverityMinor,
		},
		{
			Description: "PASSWORD",
			Regexp:      regexp.MustCompile("(?im)['\"]?" + PasswordPrefixPattern + "['\"]?[[:blank:]]{0,20}[=:]{1,3}?[[:blank:]]{0,20}[@]?" + "['\"]?" + PasswordPattern + "((?:['\"]?(?:[;,])?)?$|[[:blank:]])"),
			Tags:        []string{"password"},
			Severity:    SeverityMajor,
		},
		{
			Description: "PASSWORD_XML",
			File:        regexp.MustCompile("(?i)(.*.xml)$"),
			Regexp:      regexp.MustCompile("(?i)<" + PasswordPrefixPattern + ">(?P<secret>.{5,256})</" + PasswordPrefixPattern + ">"),
			Tags:        []string{"password"},
			Severity:    SeverityMajor,
		},
		{
			Description: "SECRET_KEY",
			Regexp:      regexp.MustCompile("(?im)" + SecretKeyPrefixPattern + "?[[:space:]_-]?(?:key|token)[[:space:]]{0,20}[=:]{1,2}[[:space:]]{0,20}['\"]?" + PasswordPattern + "(?:[[:space:];'\",]|$)"),
			Tags:        []string{"token", "key"},
			Severity:    SeverityMajor,
		},
	}
)

// AllowList list of allowed items
type AllowList struct {
	Commits     []string
	Description string
	Files       []*regexp.Regexp
	Paths       []*regexp.Regexp
	Regexes     []*regexp.Regexp
}

// Entropy sata struct
type Entropy struct {
	MinThreshold float64
	MaxThreshold float64
	Group        int
}

// Rule data struct
type Rule struct {
	AllowList   AllowList
	Description string
	Entropies   []Entropy
	File        *regexp.Regexp
	Path        *regexp.Regexp
	Regexp      *regexp.Regexp
	Severity    Severity
	Tags        []string
}
