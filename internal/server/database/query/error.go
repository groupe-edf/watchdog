package query

import "errors"

var (
	// ErrInvalidLimitation offset or limit is not correct
	ErrInvalidLimitation = errors.New("Offset or limit is not correct")
	// ErrNeedMoreArguments need more arguments
	ErrNeedMoreArguments = errors.New("Need more sql arguments")
	// ErrNotSupportType not supported SQL type error
	ErrNotSupportType = errors.New("Not supported SQL type")
	// ErrNoTableName no table name
	ErrNoTableName = errors.New("No table indicated")
)
