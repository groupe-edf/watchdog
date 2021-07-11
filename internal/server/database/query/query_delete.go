package query

import "fmt"

// WriteDelete wirte delete query
func (query *Query) WriteDelete(writer Writer) error {
	if len(query.from) <= 0 {
		return ErrNoTableName
	}
	if _, err := fmt.Fprintf(writer, "DELETE FROM %s WHERE ", query.from); err != nil {
		return err
	}
	return query.condition.Write(writer)
}

// Delete creates a delete query
func Delete(conditions ...Condition) *Query {
	query := &Query{condition: NewCondition()}
	return query.Delete(conditions...)
}
