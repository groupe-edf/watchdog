package query

import (
	"net/url"
	"strconv"
	"strings"
)

// Parser query parser
type Parser struct {
	Query  *Query
	values url.Values
}

// Extract parts from query
func (parser *Parser) Extract(values url.Values, key string) interface{} {
	for parameterKey, parameterValues := range values {
		if parameterKey == key {
			return parameterValues
		}
	}
	return nil
}

// ParseConditions parse query conditions
func (parser *Parser) ParseConditions(values url.Values) []Condition {
	conditions := make([]Condition, 0)
	valueString := parser.Extract(values, ConditionsLabel)
	if valueString != nil {
		for _, value := range valueString.([]string) {
			slices := strings.Split(value, ",")
			condition := Condition{
				Operator: OperatorType(slices[1]),
				Value:    strings.Join(slices[2:], ","),
			}
			if strings.Contains(slices[0], ".") {
				condition.Field = strings.Join(strings.Split(slices[0], "."), "\".\"")
			} else {
				condition.Field = slices[0]
			}
			conditions = append(conditions, condition)
		}
	}
	return conditions
}

// ParseInt parse integer in query
func (parser *Parser) ParseInt(values url.Values, key string, defaultValue int) int {
	valueString := parser.Extract(values, key)
	if valueString != nil {
		if value, err := strconv.Atoi(valueString.([]string)[0]); err == nil {
			return value
		}
	}
	return defaultValue
}

func (parser *Parser) ParseSort(values url.Values) []Sort {
	items := make([]Sort, 0)
	valueString := parser.Extract(values, SortLabel)
	if valueString != nil {
		for _, value := range valueString.([]string) {
			slices := strings.Split(value, ",")
			sort := Sort{
				Direction: DirectionType(slices[1]),
				Field:     slices[0],
			}
			items = append(items, sort)
		}
	}
	return items
}

// Parse query values
func (parser *Parser) Parse() {
	parser.Query = &Query{}
	parser.Query.Conditions = parser.ParseConditions(parser.values)
	parser.Query.Limit = parser.ParseInt(parser.values, LimitLabel, 10)
	parser.Query.Offset = parser.ParseInt(parser.values, OffsetLabel, 0)
	parser.Query.Sort = parser.ParseSort(parser.values)
}

// Parse query
func Parse(values url.Values) *Query {
	parser := Parser{
		values: values,
	}
	parser.Parse()
	return parser.Query
}

func NewCondition(field string, operator OperatorType, value interface{}) Condition {
	condition := Condition{
		Operator: OperatorType(operator),
		Value:    value,
	}
	if strings.Contains(field, ".") {
		condition.Field = strings.Join(strings.Split(field, "."), "\".\"")
	} else {
		condition.Field = field
	}
	return condition
}
