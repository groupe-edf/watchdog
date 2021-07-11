package security

import "github.com/groupe-edf/watchdog/internal/models"

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
)

var (
	rules = []models.Rule{
		{
			Description: "ASYMMETRIC_PRIVATE_KEY",
			Pattern:     string(`(\-){5}BEGIN[[:blank:]]*?(RSA|OPENSSH|DSA|EC|PGP)?[[:blank:]]*?PRIVATE[[:blank:]]KEY[[:blank:]]*?(BLOCK)?(\-){5}.*`),
			Tags:        []string{"key"},
			Severity:    models.SeverityMajor,
		},
		{
			Description: "AWS_ACCESS_KEY",
			Pattern:     "(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}",
			Tags:        []string{"aws"},
			Severity:    models.SeverityMajor,
		},
		{
			Description: "AWS_SECRET_KEY",
			Pattern:     string(`(?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]`),
			Tags:        []string{"aws"},
			Severity:    models.SeverityBlocker,
		},
		{
			Description: "BASE_64",
			Pattern:     "(?i)['\"]?((?:[_]?auth(?:Token|orization:[[:blank:]]Basic)?)['\"]?[[:blank:]=:]{1})[[:blank:]]*['\"]?" + Base64Pattern + "['\"]?",
			Tags:        []string{"authentication", "base64"},
			Severity:    models.SeverityMinor,
		},
		{
			Description: "CONFIDENTIAL",
			Pattern:     "(?i)CONFIDENTIAL",
			Severity:    models.SeverityInfo,
		},
		{
			Description: "CONNECTION_STRING",
			Pattern:     "(?i)" + ConnectionString,
			Severity:    models.SeverityMajor,
		},
		{
			Description: "ENTROPY",
			File:        "(?i).*.sql$",
			Entropies: []models.Entropy{
				{
					MaxThreshold: 8.0,
					MinThreshold: 4.0,
				},
			},
			Pattern:  string(`[0-9a-zA-Z-_!{}$.\/=]{8,120}`),
			Tags:     []string{"entropy"},
			Severity: models.SeverityInfo,
		},
		{
			Description: "HTPASSWD",
			File:        "(?i).htpasswd$",
			Pattern:     "(?i)[0-9a-zA-Z-_!{}$.=]{4,120}:" + PasswordPattern,
			Severity:    models.SeverityMinor,
		},
		{
			Description: "PASSWORD",
			Pattern:     "(?im)['\"]?" + PasswordPrefixPattern + "['\"]?[[:blank:]]{0,20}[=:]{1,3}?[[:blank:]]{0,20}[@]?" + "['\"]?" + PasswordPattern + "((?:['\"]?(?:[;,])?)?$|[[:blank:]])",
			Tags:        []string{"password"},
			Severity:    models.SeverityMajor,
		},
		{
			Description: "PASSWORD_XML",
			File:        "(?i)(.*.xml)$",
			Pattern:     "(?i)<" + PasswordPrefixPattern + ">(?P<secret>.{5,256})</" + PasswordPrefixPattern + ">",
			Tags:        []string{"password"},
			Severity:    models.SeverityMajor,
		},
		{
			Description: "SECRET_KEY",
			Pattern:     "(?im)" + SecretKeyPrefixPattern + "?[[:space:]_-]?(?:key|token)[[:space:]]{0,20}[=:]{1,2}[[:space:]]{0,20}['\"]?" + PasswordPattern + "(?:[[:space:];'\",]|$)",
			Tags:        []string{"token", "key"},
			Severity:    models.SeverityMajor,
		},
	}
)
