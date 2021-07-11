package query

// OrCondition and condition
type OrCondition []Condition

// And condition
func (condition OrCondition) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition OrCondition) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition OrCondition) IsValid() bool {
	return len(condition) > 0
}

func (condition OrCondition) Write(writer Writer) error {
	return nil
}

// Or condition
func Or(conditions ...Condition) Condition {
	var result = make(OrCondition, 0, len(conditions))
	for _, condition := range conditions {
		if condition == nil || !condition.IsValid() {
			continue
		}
		result = append(result, condition)
	}
	return result
}
