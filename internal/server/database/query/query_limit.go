package query

import "fmt"

// Limit data struct
type Limit struct {
	limit  int
	offset int
}

// WriteLimit wirte limit in sql query
func (query *Query) WriteLimit(writer Writer) error {
	if query.limit != nil {
		limit := query.limit
		if limit.offset < 0 || limit.limit <= 0 {
			return ErrInvalidLimitation
		}
		query.limit = nil
		defer func() {
			query.limit = limit
		}()
		limitWriter := writer.(*BytesWriter)
		if limit.offset == 0 {
			fmt.Fprint(limitWriter, " LIMIT ", limit.limit)
		} else {
			fmt.Fprintf(limitWriter, " LIMIT %v OFFSET %v", limit.limit, limit.offset)
		}
	}
	return nil
}
