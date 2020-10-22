package security

import "regexp"

// Scanner scanner interface
type Scanner interface {
	Scan(content string) []string
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}

// BaseScanner data struct
type BaseScanner struct {
	Provider string
	Regexp   *regexp.Regexp
	Matches  []string
}
