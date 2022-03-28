package query

import (
	"net/url"
	"strconv"
)

const (
	// ConditionsLabel label
	ConditionsLabel = "conditions"
	// LimitLabel label
	LimitLabel = "limit"
	// OffsetLabel label
	OffsetLabel = "offset"
	// FullTextQueryLabel label
	FullTextQueryLabel = "query"
	// SearchLabel label
	SearchLabel = "search"
	// SortLabel label
	SortLabel = "sort"
)

// Fetcher repository fetcher interface
type Fetcher interface {
	Count() int
	Find(q *Query) ([]interface{}, error)
}

// Query data struct
type Query struct {
	Conditions []Condition `url:"conditions,omitempty" json:"conditions,omitempty"`
	Limit      int         `url:"limit,omitempty" json:"limit,omitempty"`
	Offset     int         `url:"offset,omitempty" json:"offset,omitempty"`
	Search     string      `url:"search,omitempty" json:"search,omitempty"`
	Sort       []Sort      `url:"sort,omitempty" json:"sort,omitempty"`
}

// AddCondition add filter condition
func (query *Query) AddCondition(condition Condition) {
	query.Conditions = append(query.Conditions, condition)
}

// AddSort add sorter
func (query *Query) AddSort(sort Sort) {
	query.Sort = append(query.Sort, sort)
}

// SetSearch set search query
func (query *Query) SetSearch(search string) {
	query.Search = search
}

// Result query result
type Result struct {
	Items      []interface{}
	Limit      int
	Offset     int
	TotalItems int
}

// ParseQuery parse query string and retur Query
func ParseQuery(queryString string) *Query {
	query := &Query{
		Limit:  10,
		Offset: 0,
	}
	rawQuery, err := url.Parse(queryString)
	if err != nil {
		return query
	}
	parsedQuery := rawQuery.Query()
	if limitString := parsedQuery.Get(LimitLabel); limitString != "" {
		if limit, err := strconv.Atoi(limitString); err == nil {
			query.Limit = limit
		}
	}
	if offsetString := parsedQuery.Get(OffsetLabel); offsetString != "" {
		if offset, err := strconv.Atoi(offsetString); err == nil {
			query.Offset = offset
		}
	}
	return query
}
