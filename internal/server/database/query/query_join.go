package query

// Join data struct
type Join struct {
	joinTable     interface{}
	joinType      string
	joinCondition Condition
}

// Join sets join table and conditions
func (query *Query) Join(joinType string, joinTable, condition interface{}) *Query {
	switch condition.(type) {
	case Condition:
		query.joins = append(query.joins, Join{joinTable, joinType, condition.(Condition)})
	case string:
		query.joins = append(query.joins, Join{joinTable, joinType, Expression(condition.(string))})
	}
	return query
}
