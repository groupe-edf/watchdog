package query

import (
	"fmt"
	"sort"
)

// Equal condition
type Equal map[string]interface{}

var _ Condition = Equal{}

// And condition
func (condition Equal) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition Equal) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition Equal) IsValid() bool {
	return len(condition) > 0
}

// SortKeys sort condition keys
func (condition Equal) SortKeys() []string {
	keys := make([]string, 0, len(condition))
	for key := range condition {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Write query
func (condition Equal) Write(writer Writer) error {
	var i = 0
	for _, key := range condition.SortKeys() {
		value := condition[key]
		switch value.(type) {
		default:
			if _, err := fmt.Fprintf(writer, "%s = ?", key); err != nil {
				return err
			}
			writer.Append(value)
		}
		if i != len(condition)-1 {
			if _, err := fmt.Fprint(writer, " AND "); err != nil {
				return err
			}
		}
		i = i + 1
	}
	return nil
}
