package query

import "fmt"

// AndCondition and condition
type AndCondition []Condition

// And condition
func (condition AndCondition) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition AndCondition) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition AndCondition) IsValid() bool {
	return len(condition) > 0
}

func (condition AndCondition) Write(writer Writer) error {
	for i, conditionValue := range condition {
		_, isOrCondition := conditionValue.(OrCondition)
		_, isExpression := conditionValue.(ExpressionCondition)
		wrap := isOrCondition || isExpression
		if wrap {
			fmt.Fprint(writer, "(")
		}
		err := conditionValue.Write(writer)
		if err != nil {
			return err
		}
		if wrap {
			fmt.Fprint(writer, ")")
		}
		if i != len(condition)-1 {
			fmt.Fprint(writer, " AND ")
		}
	}
	return nil
}

var _ Condition = AndCondition{}

// And condition
func And(conditions ...Condition) Condition {
	var result = make(AndCondition, 0, len(conditions))
	for _, condition := range conditions {
		if condition == nil || !condition.IsValid() {
			continue
		}
		result = append(result, condition)
	}
	return result
}
