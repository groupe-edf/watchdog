package security

import (
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Options scanner options
type Options struct {
	AllowList AllowList
}

// Scanner scanner interface
type Scanner interface {
	AddAllowedFiles(files *regexp.Regexp)
	Scan(commit *object.Commit) (leaks []Leak, err error)
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}
