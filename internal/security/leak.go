package security

// Leak data struct
type Leak struct {
	File        string   `json:"file"`
	Line        string   `json:"-"`
	LineNumber  int      `json:"line_number"`
	Offender    string   `json:"-"`
	Remediation string   `json:"remediation,omitempty"`
	Rule        string   `json:"rule"`
	Severity    Severity `json:"severity"`
	Tags        []string `json:"tags,omitempty"`
}
