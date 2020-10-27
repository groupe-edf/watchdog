package logging

// Options logger options
type Options struct {
	LogsFormat string
	LogsLevel  string
	LogsPath   string
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
func New(options Options) Interface {
	return NewLogrusLogger(options)
}
