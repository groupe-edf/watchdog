package query

const (
	// OffsetType pagination offset type
	OffsetType = "offset"
	// CursorType pagination cursor type
	CursorType = "cursor"
)

// Pagination data struct
type Pagination struct {
	Limit      int
	Offset     int
	TotalItems int
}
