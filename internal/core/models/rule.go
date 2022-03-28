package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

type Fragment struct {
	CommitSHA      string
	FilePath       string
	Keywords       map[string]bool
	NewlineIndices [][]int
	Raw            string
}

type Location struct {
	startLine      int
	endLine        int
	startColumn    int
	endColumn      int
	startLineIndex int
	endLineIndex   int
}

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

func (rule Rule) Scan(fragment Fragment) (leaks []Leak) {
	// fix: compile regexp once at start
	pattern := regexp.MustCompile(rule.Pattern)
	matchIndices := pattern.FindAllStringIndex(fragment.Raw, -1)
	for _, matchIndex := range matchIndices {
		// extract secret from match
		offender := strings.Trim(fragment.Raw[matchIndex[0]:matchIndex[1]], "\n")
		location := rule.location(fragment, matchIndex)
		if matchIndex[1] > location.endLineIndex {
			location.endLineIndex = matchIndex[1]
		}
		leak := Leak{
			CommitHash: fragment.CommitSHA,
			File:       fragment.FilePath,
			LineNumber: location.startColumn,
			Location:   location,
			Offender:   offender,
			Rule:       rule,
			Tags:       rule.Tags,
			Severity:   rule.Severity,
			SecretHash: fmt.Sprintf("%s:%s:%s:%d:%d", time.Now(), fragment.CommitSHA, fragment.FilePath, location.startLine, rule.ID),
		}
		leaks = append(leaks, leak)
	}
	return leaks
}

func (rule *Rule) Validate() error {
	if rule.Pattern == "" {
		return errors.New("pattern is mandatory")
	}
	return nil
}

func (rule *Rule) location(fragment Fragment, matchIndex []int) Location {
	var (
		prevNewLine int
		location    Location
		lineSet     bool
		_lineNum    int
	)
	start := matchIndex[0]
	end := matchIndex[1]
	location.startLineIndex = 0
	for lineNum, pair := range fragment.NewlineIndices {
		_lineNum = lineNum
		newLineByteIndex := pair[0]
		if prevNewLine <= start && start < newLineByteIndex {
			lineSet = true
			location.startLine = lineNum
			location.endLine = lineNum
			location.startColumn = (start - prevNewLine) + 1 // +1 because counting starts at 1
			location.startLineIndex = prevNewLine
			location.endLineIndex = newLineByteIndex
		}
		if prevNewLine < end && end <= newLineByteIndex {
			location.endLine = lineNum
			location.endColumn = (end - prevNewLine)
			location.endLineIndex = newLineByteIndex
		}
		prevNewLine = pair[0]
	}
	if !lineSet {
		location.startColumn = (start - prevNewLine) + 1 // +1 because counting starts at 1
		location.endColumn = (end - prevNewLine)
		location.startLine = _lineNum + 1
		location.endLine = _lineNum + 1
		i := 0
		for end+i < len(fragment.Raw) {
			if fragment.Raw[end+i] == '\n' {
				break
			}
			if fragment.Raw[end+i] == '\r' {
				break
			}
			i++
		}
		location.endLineIndex = end + i
	}
	return location
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
