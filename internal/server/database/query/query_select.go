package query

import "fmt"

// WriteSelect write select statement
func (query *Query) WriteSelect(writer Writer) error {
	if len(query.from) <= 0 && !query.nested {
		return ErrNoTableName
	}
	if _, err := fmt.Fprint(writer, "SELECT "); err != nil {
		return err
	}
	if len(query.columns) > 0 {
		for index, column := range query.columns {
			if _, err := fmt.Fprint(writer, column); err != nil {
				return err
			}
			if index != len(query.columns)-1 {
				if _, err := fmt.Fprint(writer, ", "); err != nil {
					return err
				}
			}
		}
	} else {
		if _, err := fmt.Fprint(writer, "*"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(writer, " FROM ", query.from); err != nil {
		return err
	}
	for _, join := range query.joins {
		joinQuery, ok := join.joinTable.(*Query)
		if ok {
			if _, err := fmt.Fprintf(writer, " %s JOIN (", join.joinType); err != nil {
				return err
			}
			if err := joinQuery.Write(writer); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(writer, ") ON "); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(writer, " %s JOIN %s ON ", join.joinType, join.joinTable); err != nil {
				return err
			}
		}
		if err := join.joinCondition.Write(writer); err != nil {
			return err
		}
	}
	if query.condition.IsValid() {
		if _, err := fmt.Fprint(writer, " WHERE "); err != nil {
			return err
		}

		if err := query.condition.Write(writer); err != nil {
			return err
		}
	}
	if len(query.order) > 0 {
		if _, err := fmt.Fprint(writer, " ORDER BY ", query.order); err != nil {
			return err
		}
	}
	if query.limit != nil {
		if err := query.WriteLimit(writer); err != nil {
			return err
		}
	}
	return nil
}
