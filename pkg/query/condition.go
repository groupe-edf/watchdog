package query

// OperatorType operator type
type OperatorType string

const (
	// Between operator
	Between OperatorType = "between"
	// Empty operator
	Empty = "empty"
	// Equal operator
	Equal = "eq"
	// Exists operator
	Exists = "exists"
	// GreaterThan operator
	GreaterThan = "gt"
	// In operator
	In = "in"
	// Like operator
	Like = "like"
	// LowerThan operator
	LowerThan = "lt"
	// NotEqual operator
	NotEqual = "ne"
	// NotIn operator
	NotIn = "nin"
)

var (
	// Operators list of query available operators
	Operators = []OperatorType{
		Between,
		Empty,
		Equal,
		Exists,
		GreaterThan,
		In,
		Like,
		LowerThan,
		NotEqual,
		NotIn,
	}
)

// Condition query condition
type Condition struct {
	Field    string
	Operator OperatorType
	Value    interface{}
}
