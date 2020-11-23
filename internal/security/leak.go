package security

// Leak data struct
type Leak struct {
	File        string   `json:"file"`
	Line        string   `json:"line"`
	LineNumber  int      `json:"line_number"`
	Offender    string   `json:"offender"`
	Remediation string   `json:"remediation"`
	Rule        string   `json:"rule"`
	Severity    Severity `json:"severity"`
	Tags        []string `json:"tags,omitempty"`
}
