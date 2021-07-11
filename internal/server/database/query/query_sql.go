package query

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ToBoundSQL return sql bounded with arguments
func (query *Query) ToBoundSQL() (string, error) {
	writer := NewWriter()
	if err := query.Write(writer); err != nil {
		return "", err
	}
	return ConvertToBoundSQL(writer.String(), writer.args)
}

// ConvertToBoundSQL bound arguments with values
func ConvertToBoundSQL(sql string, args []interface{}) (string, error) {
	buffer := strings.Builder{}
	var i, j, start int
	for ; i < len(sql); i++ {
		if sql[i] == '?' {
			_, err := buffer.WriteString(sql[start:i])
			if err != nil {
				return "", err
			}
			start = i + 1
			if len(args) == j {
				return "", ErrNeedMoreArguments
			}
			arg := args[j]
			if NoNeedQuote(arg) {
				_, err = fmt.Fprint(&buffer, arg)
			} else {
				_, err = fmt.Fprintf(&buffer, "'%v'", strings.Replace(fmt.Sprintf("%v", arg), "'", "''", -1))
			}
			if err != nil {
				return "", err
			}
			j = j + 1
		}
	}
	_, err := buffer.WriteString(sql[start:])
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// NoNeedQuote check if value do not need quote
func NoNeedQuote(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	case bool:
		return true
	case string:
		return false
	case time.Time, *time.Time:
		return false
	}
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Bool:
		return true
	case reflect.String:
		return false
	}
	return false
}
