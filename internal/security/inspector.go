package security

import (
	"path/filepath"
	"strings"
)

const (
	// IsPositive string is secret
	IsPositive int = iota
	// IsFile string is a path
	IsFile
	// IsFunction string is function
	IsFunction
	// IsPlaceholder string is placeholder
	IsPlaceholder
	// IsVariable string is variable
	IsVariable
	// PerCharThreshold entropy per character threshold
	PerCharThreshold = 3
)

var (
	// SupportedLanguages list of supported languages
	SupportedLanguages = []string{"go", "groovy", "java", "js", "py"}
)

type AfterInspector interface {
	After(filePath string, line string, secret string) int
}

type BeforeInspector interface {
	Before()
}

type Inspector interface {
	AfterInspector
	BeforeInspector
}

// IsFalsePositive check if secret is a false positive
func IsFalsePositive(filePath string, line string, secret string) int {
	// Secret is a variable
	if strings.HasPrefix(secret, "$") && !strings.Contains(secret[2:], "$") {
		return IsVariable
	}
	// Secret is a placeholder
	if strings.Contains(secret, "{{") || strings.Contains(secret, "}}") {
		return IsPlaceholder
	}
	// Secret is a placeholder
	if strings.HasPrefix(secret, "{") && strings.HasSuffix(secret, "}") {
		if len(secret) < 32 {
			return IsPlaceholder
		}
	}
	// Secret is a placeholder
	if strings.HasPrefix(secret, "${") && strings.HasSuffix(secret, "}") {
		return IsPlaceholder
	}
	extension := filepath.Ext(filePath)
	openBrackets := strings.Count(secret, "(")
	closeBrackets := strings.Count(secret, ")")
	// Secret is method or function
	if IsSupportedLanguage(extension) {
		if openBrackets >= 1 {
			if openBrackets == closeBrackets {
				return IsFunction
			}
		}
		if strings.Count(line, "\""+secret+"\"") < 1 {
			return IsVariable
		}
	}
	if strings.HasSuffix(secret, ";") {
		return IsVariable
	}
	return IsPositive
}

// IsSupportedLanguage check if extension is suported
func IsSupportedLanguage(language string) bool {
	for _, supported := range SupportedLanguages {
		if language == "."+supported {
			return true
		}
	}
	return false
}
