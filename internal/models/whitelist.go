package models

// Whitelist list of allowed items
type Whitelist struct {
	ID           int64  `json:"id,omitempty"`
	Commits      List   `json:"commits,omitempty"`
	Description  string `json:"description,omitempty"`
	Files        List   `json:"files,omitempty"`
	Paths        List   `json:"paths,omitempty"`
	Regexes      List   `json:"regexes,omitempty"`
	Repositories List   `json:"repositories,omitempty"`
}
