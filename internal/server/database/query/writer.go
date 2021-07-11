package query

import (
	"io"
	"strings"
)

// Writer query writer interface
type Writer interface {
	io.Writer
	Append(...interface{})
}

var _ Writer = NewWriter()

// BytesWriter default writer
type BytesWriter struct {
	*strings.Builder
	args []interface{}
}

// NewWriter creates a new string writer
func NewWriter() *BytesWriter {
	writer := &BytesWriter{
		Builder: &strings.Builder{},
	}
	return writer
}

// Append appends args to Writer
func (writer *BytesWriter) Append(args ...interface{}) {
	writer.args = append(writer.args, args...)
}

// Args returns args
func (writer *BytesWriter) Args() []interface{} {
	return writer.args
}
