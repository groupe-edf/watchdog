package query

// DirectionType sort direction
type DirectionType string

const (
	// Ascending sort
	Ascending DirectionType = "asc"
	// Descending sort
	Descending = "desc"
)

// Sort query sort
type Sort struct {
	Direction DirectionType
	Field     string
}
