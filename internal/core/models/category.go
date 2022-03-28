package models

// Category entity
type Category struct {
	ID          *int64 `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Extension   string `json:"extension,omitempty"`
	Left        int64  `json:"left,omitempty"`
	Level       int    `json:"level,omitempty"`
	ParentID    int64  `json:"parent_id,omitempty"`
	Right       int64  `json:"right,omitempty"`
	Title       string `json:"title,omitempty"`
	Value       string `json:"value,omitempty"`
}
