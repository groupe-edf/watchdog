package security

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
)

// Scanner scanner interface
type Scanner interface {
	Scan(gragment models.Fragment) (leaks []models.Leak, err error)
	SetWhitelist(whitelist models.Whitelist)
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}
