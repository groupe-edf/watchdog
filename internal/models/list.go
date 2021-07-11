package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type List []string

func (list *List) Append(items string) {
	*list = append(*list, items)
}

// Scan implements sql.Scanner so list can be read from databases transparently
func (list *List) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		return nil
	case string:
		if src == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(src), &list); err != nil {
			return err
		}
	}
	return nil
}

func (list List) String() string {
	return strings.Join(list, ",")
}

func (list *List) Value() (driver.Value, error) {
	return list.String(), nil
}

func (list *List) MarshalJSON() ([]byte, error) {
	return json.Marshal((*[]string)(list))
}

func (list *List) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*[]string)(list)); err != nil {
		return err
	}
	return nil
}
