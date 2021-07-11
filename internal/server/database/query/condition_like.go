package query

import "fmt"

// Like and condition
type Like [2]string

var _ Condition = Like{"", ""}

// And condition
func (condition Like) And(conditions ...Condition) Condition {
	return And(condition, And(conditions...))
}

// Or condition
func (condition Like) Or(conditions ...Condition) Condition {
	return Or(condition, Or(conditions...))
}

// IsValid check if condition is valid
func (condition Like) IsValid() bool {
	return len(condition) > 0
}

func (condition Like) Write(writer Writer) error {
	if _, err := fmt.Fprintf(writer, "%s ILIKE ?", condition[0]); err != nil {
		return err
	}
	if condition[1][0] == '%' || condition[1][len(condition[1])-1] == '%' {
		writer.Append(condition[1])
	} else {
		writer.Append("%" + condition[1] + "%")
	}
	return nil
}
