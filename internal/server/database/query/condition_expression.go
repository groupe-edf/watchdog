package query

import "fmt"

// ExpressionCondition sql expression
type ExpressionCondition struct {
	sql  string
	args []interface{}
}

var _ Condition = ExpressionCondition{}

// Expression create an expression condition
func Expression(sql string, args ...interface{}) Condition {
	return ExpressionCondition{sql, args}
}

// And condition
func (condition ExpressionCondition) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition ExpressionCondition) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition ExpressionCondition) IsValid() bool {
	return len(condition.sql) > 0
}

func (condition ExpressionCondition) Write(writer Writer) error {
	if _, err := fmt.Fprint(writer, condition.sql); err != nil {
		return err
	}
	writer.Append(condition.args...)
	return nil
}
