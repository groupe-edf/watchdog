package security

import (
	"regexp"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

// IsFalsePositiveFunc function to check false positive secrets
type IsFalsePositiveFunc func(file string, line string, secret string) int

// RegexScanner data struct
type RegexScanner struct {
	Logger    logging.Interface
	Rules     []models.Rule
	Whitelist models.Whitelist
}

func (scanner *RegexScanner) Scan(fragment models.Fragment) (leaks []models.Leak, err error) {
	fragment.Keywords = make(map[string]bool)
	fragment.NewlineIndices = regexp.MustCompile("\n").FindAllStringIndex(fragment.Raw, -1)
	for _, rule := range scanner.Rules {
		leaks = append(leaks, rule.Scan(fragment)...)
	}
	return leaks, err
}

func (scanner *RegexScanner) SetWhitelist(whitelist models.Whitelist) {
	scanner.Whitelist = whitelist
}

// NewRegexScanner create new regular expression
func NewRegexScanner(logger logging.Interface, rules []models.Rule, whitelist models.Whitelist) *RegexScanner {
	return &RegexScanner{
		Logger:    logger,
		Rules:     rules,
		Whitelist: whitelist,
	}
}
