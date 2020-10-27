package security

// Options scanner options
type Options struct {
	AllowList AllowList
}

// Scanner scanner interface
type Scanner interface {
	Scan(content string) []string
}

// VerifiedScanner verify secret
type VerifiedScanner interface {
	Verify() bool
}
