package security

import (
	"github.com/groupe-edf/watchdog/internal/models"
)

// Scanner scanner interface
type Scanner interface {
	Scan(commit *models.Commit) (leaks []models.Leak, err error)
	SetWhitelist(whitelist models.Whitelist)
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}
