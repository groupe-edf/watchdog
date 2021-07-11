package query

// Condition interface
type Condition interface {
	Write(Writer) error
	And(...Condition) Condition
	Or(...Condition) Condition
	IsValid() bool
}

// EmptyCondition data struct
type EmptyCondition struct{}

// And condition
func (EmptyCondition) And(conditions ...Condition) Condition {
	return And(conditions...)
}

// IsValid check if condition is valid
func (EmptyCondition) IsValid() bool {
	return false
}

// Or condition
func (EmptyCondition) Or(conditions ...Condition) Condition {
	return Or(conditions...)
}

// Write query
func (EmptyCondition) Write(writer Writer) error {
	return nil
}

// NewCondition return empty condition
func NewCondition() Condition {
	return EmptyCondition{}
}
