package query

import (
	"fmt"
	"strings"
)

// In condition
type InCondition struct {
	columns string
	values  []interface{}
}

var _ Condition = InCondition{}

// In generates IN condition
func In(columns string, values ...interface{}) Condition {
	return InCondition{columns, values}
}

// And condition
func (condition InCondition) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

func (condition InCondition) HandleBlank(w Writer) error {
	_, err := fmt.Fprint(w, "0=1")
	return err
}

func (condition InCondition) IsValid() bool {
	return len(condition.columns) > 0 && len(condition.values) > 0
}

// Or condition
func (condition InCondition) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

func (condition InCondition) Write(w Writer) error {
	if len(condition.values) <= 0 {
		return condition.HandleBlank(w)
	}
	switch condition.values[0].(type) {
	case []string:
		values := condition.values[0].([]string)
		if len(values) <= 0 {
			return condition.HandleBlank(w)
		}
		questionMark := strings.Repeat("?,", len(values))
		if _, err := fmt.Fprintf(w, "%s IN (%s)", condition.columns, questionMark[:len(questionMark)-1]); err != nil {
			return err
		}
		for _, value := range values {
			w.Append(value)
		}
	}
	return nil
}
