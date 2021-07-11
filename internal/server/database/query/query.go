package query

import (
	"database/sql"
	"fmt"
	"strings"

	routeQuery "github.com/groupe-edf/watchdog/internal/server/query"
)

// OperationType sql operation type
type OperationType byte

const (
	// ConditionType condition operation
	ConditionType OperationType = iota
	// SelectType select operation
	SelectType
	// InsertType insert operation
	InsertType
	// UpdateType update operation
	UpdateType
	// DeleteType delete operation
	DeleteType
	// SetOperationType set operation type
	SetOperationType
)

// Query provides all the chainable relational query builder methods
type Query struct {
	columns   []string
	condition Condition
	from      string
	group     string
	having    string
	joins     []Join
	limit     *Limit
	nested    bool
	offset    int
	operation OperationType
	order     string
	sql       string
	store     *sql.DB
	table     string
}

// Delete query
func (query *Query) Delete(conditions ...Condition) *Query {
	query.condition = query.condition.And(conditions...)
	query.operation = DeleteType
	return query
}

// From set query from clause
func (query *Query) From(from interface{}, alias ...string) *Query {
	switch from.(type) {
	case string:
		query.from = from.(string)
		if len(alias) > 0 {
			query.from = query.from + " " + alias[0]
		}
	}
	return query
}

// Limit set query limit
func (query *Query) Limit(limit int, offset ...int) *Query {
	query.limit = &Limit{limit: limit}
	if len(offset) > 0 {
		query.limit.offset = offset[0]
	}
	return query
}

// Offset set query offset
func (query *Query) Offset(offset int) *Query {
	query.offset = offset
	return query
}

// OrderBy set query order by clause
func (query *Query) OrderBy(order string) *Query {
	query.order = order
	return query
}

// Select set query select columns
func (query *Query) Select(columns ...string) *Query {
	query.columns = append(query.columns, columns...)
	if query.operation == ConditionType {
		query.operation = SelectType
	}
	return query
}

// SQL set sql query
func (query *Query) SQL(sql string) *Query {
	query.sql = sql
	return query
}

// ToSQL return sql query
func (query *Query) ToSQL() (string, []interface{}, error) {
	writer := NewWriter()
	if err := query.Write(writer); err != nil {
		return "", nil, err
	}
	var sql = writer.String()
	return sql, writer.args, nil
}

// Where set where conditions
func (query *Query) Where(condition Condition) *Query {
	if query.condition.IsValid() {
		query.condition = query.condition.And(condition)
	} else {
		query.condition = condition
	}
	return query
}

func (query *Query) WithRouteQuery(q *routeQuery.Query) *Query {
	if len(q.Conditions) > 0 {
		for _, condition := range q.Conditions {
			switch condition.Operator {
			case routeQuery.Equal:
				query.Where(Equal{fmt.Sprintf(`%v`, condition.Field): condition.Value})
			case routeQuery.In:
				slices := strings.Split(condition.Value.(string), ",")
				query.Where(In(fmt.Sprintf(`%v`, condition.Field), slices))
			case routeQuery.Like:
				query.Where(Like{fmt.Sprintf(`%v`, condition.Field), condition.Value.(string)})
			}
		}
	}
	if len(q.Sort) > 0 {
		query.OrderBy(fmt.Sprintf(`"%s" %s`, q.Sort[0].Field, q.Sort[0].Direction))
	}
	return query
}

// WithStore set store
func (query *Query) WithStore(store *sql.DB) *Query {
	query.store = store
	return query
}

func (query *Query) Write(writer Writer) error {
	switch query.operation {
	case SelectType:
		return query.WriteSelect(writer)
	case DeleteType:
		return query.WriteDelete(writer)
	}
	return ErrNotSupportType
}

// NewQuery create new sql query
func NewQuery(from string) *Query {
	query := &Query{}
	query.From(from)
	return query
}

// Select creates a select Builder
func Select(columns ...string) *Query {
	query := &Query{condition: NewCondition()}
	return query.Select(columns...)
}
