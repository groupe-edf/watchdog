package models

import (
	"bytes"
	"database/sql/driver"
)

// Score type used by severity and confidence values
type Score int

//go:generate stringer -type=Score
const (
	// SeverityLow severity or confidence
	SeverityLow Score = iota
	// SeverityMedium severity or confidence
	SeverityMedium
	// SeverityHigh severity or confidence
	SeverityHigh
)

// MarshalJSON marshal score to json
func (score Score) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(score.String())
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (score *Score) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*score = ParseScore(value)
	return nil
}

func (score Score) String() string {
	switch score {
	case SeverityHigh:
		return "HIGHT"
	case SeverityMedium:
		return "MEDIUM"
	case SeverityLow:
		return "LOW"
	}
	return ""
}

func (score *Score) Value() (driver.Value, error) {
	return score.String(), nil
}

// UnmarshalJSON unmarshal json to score
func (score *Score) UnmarshalJSON(raw []byte) error {
	runes := bytes.Runes(raw)
	if runes[0] == '"' && runes[len(runes)-1] == '"' {
		runes = runes[1 : len(runes)-1]
	}
	*score = ParseScore(string(runes))
	return nil
}

// ParseScore parse score from string input
func ParseScore(score interface{}) Score {
	switch score {
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	}
	return SeverityLow
}
