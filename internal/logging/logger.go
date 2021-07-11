package logging

import "io"

// Options logger options
type Options struct {
	LogsFormat       string
	LogsLevel        string
	LogsOutput       io.Writer
	LogsPath         string
	LogsReportCaller bool
}

// Interface logger interface
type Interface interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	WithField(key string, value interface{}) Interface
	WithFields(fields Fields) Interface
}

// Fields convenience type for adding multiple fields to a log statement.
type Fields map[string]interface{}

// New return default logger
func New(options Options) *LogrusLogger {
	return NewLogrusLogger(options)
}
