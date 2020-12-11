package issue

import "bytes"

const (
	// SeverityLow severity or confidence
	SeverityLow Score = iota
	// SeverityMedium severity or confidence
	SeverityMedium
	// SeverityHigh severity or confidence
	SeverityHigh
)

// Score type used by severity and confidence values
type Score int

// MarshalJSON return score as json
func (score Score) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(score.String())
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (score Score) String() string {
	switch score {
	case SeverityHigh:
		return "high"
	case SeverityMedium:
		return "medium"
	case SeverityLow:
		return "low"
	}
	return "undefined"
}

// ParseScore parse score from string input
func ParseScore(score string) Score {
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
