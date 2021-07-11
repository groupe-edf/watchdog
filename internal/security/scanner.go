package security

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/groupe-edf/watchdog/internal/models"
)

// Scanner scanner interface
type Scanner interface {
	Scan(commit *object.Commit) (leaks []models.Leak, err error)
	SetWhitelist(whitelist models.Whitelist)
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}
