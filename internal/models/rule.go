package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

const (
	SeverityBlocker  Severity = "BLOCKER"
	SeverityCritical Severity = "CRITICAL"
	SeverityInfo     Severity = "INFO"
	SeverityMajor    Severity = "MAJOR"
	SeverityMinor    Severity = "MINOR"
)

// Severity rule severity
type Severity string

// Entropy sata struct
type Entropy struct {
	MinThreshold float64
	MaxThreshold float64
	Group        int
}

// Rule data struct
type Rule struct {
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	CreatedBy   *uuid.UUID     `json:"created_by,omitempty"`
	Description string         `json:"description,omitempty"`
	DisplayName string         `json:"display_name,omitempty"`
	Enabled     bool           `json:"enabled,omitempty"`
	Entropies   []Entropy      `json:"entropies,omitempty"`
	File        string         `json:"file,omitempty"`
	ID          int64          `json:"id"`
	Name        string         `json:"name,omitempty"`
	Path        *regexp.Regexp `json:"path,omitempty"`
	Pattern     string         `json:"pattern,omitempty"`
	Severity    Severity       `json:"severity,omitempty"`
	Tags        List           `json:"tags,omitempty"`
	Whitelist   Whitelist      `json:"whitelist,omitempty"`
}

type ScanEntry struct {
	Groups   []string
	Line     string
	Offender string
}

func (rule *Rule) Scan(content string) (entries []ScanEntry) {
	pattern := regexp.MustCompile(rule.Pattern)
	matches := pattern.FindAllIndex([]byte(content), -1)
	if len(matches) != 0 {
		for _, match := range matches {
			line := extractLine(match[0], match[1], content)
			offender := content[match[0]:match[1]]
			groups := pattern.FindStringSubmatch(offender)
			names := pattern.SubexpNames()
			for index, group := range groups {
				if index != 0 && names[index] == "secret" {
					offender = group
					break
				}
			}
			entries = append(entries, ScanEntry{
				Groups:   groups,
				Line:     line,
				Offender: offender,
			})
		}
	}
	return entries
}

func (rule *Rule) Validate() error {
	if rule.Pattern == "" {
		return errors.New("pattern is mandatory")
	}
	return nil
}

func extractLine(start int, end int, content string) string {
	for start != 0 && content[start] != '\n' {
		start--
	}
	if content[start] == '\n' {
		start++
	}
	for end < len(content)-1 && content[end] != '\n' {
		end++
	}
	return content[start:end]
}

func NewRule(description string, file string, pattern string, severity string, tags []string) *Rule {
	rule := &Rule{
		Description: description,
		Pattern:     pattern,
		Severity:    Severity(severity),
		Tags:        tags,
	}
	if file != "" {
		rule.File = file
	}
	if err := rule.Validate(); err != nil {
		return nil
	}
	return rule
}
