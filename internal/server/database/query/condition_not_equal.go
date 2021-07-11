package query

import (
	"fmt"
	"sort"
)

// NotEqual defines not equal conditions
type NotEqual map[string]interface{}

var _ Condition = NotEqual{}

// And condition
func (condition NotEqual) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition NotEqual) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition NotEqual) IsValid() bool {
	return len(condition) > 0
}

// SortKeys sort condition keys
func (condition NotEqual) SortKeys() []string {
	keys := make([]string, 0, len(condition))
	for key := range condition {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Write query
func (condition NotEqual) Write(writer Writer) error {
	var args = make([]interface{}, 0, len(condition))
	var i = 0
	for _, key := range condition.SortKeys() {
		value := condition[key]
		switch value.(type) {
		default:
			if _, err := fmt.Fprintf(writer, "%s<>?", key); err != nil {
				return err
			}
			args = append(args, value)
		}
		if i != len(condition)-1 {
			if _, err := fmt.Fprint(writer, " AND "); err != nil {
				return err
			}
		}
		i = i + 1
	}
	writer.Append(args...)
	return nil
}
