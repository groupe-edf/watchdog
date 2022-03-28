package logging

import "sync"

var (
	NamedLoggers loggers
)

type loggers struct {
	sync.Map
}
