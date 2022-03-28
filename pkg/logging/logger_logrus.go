package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewLogrusLogger create new logrus logger
func NewLogrusLogger(options Options) *LogrusLogger {
	logger := logrus.New()
	logLevel, _ := logrus.ParseLevel(options.LogsLevel)
	logger.SetLevel(logLevel)
	logger.SetReportCaller(options.LogsReportCaller)
	if options.LogsPath != "" {
		logFile, err := os.OpenFile(filepath.Clean(options.LogsPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			logger.Fatalf("failed to open log file: %v", err)
		}
		logger.SetOutput(logFile)
	} else {
		if options.LogsOutput != nil {
			logger.SetOutput(options.LogsOutput)
		} else {
			logger.SetOutput(io.Discard)
		}
	}
	if options.LogsFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				slices := strings.Split(f.File, "/")
				file := slices[len(slices)-1]
				return "", fmt.Sprintf("%s:%d", file, f.Line)
			},
		})
	}
	return &LogrusLogger{logger}
}

// Logrus logger adapter
type LogrusLogger struct {
	*logrus.Logger
}

// WithField append field to log entry
func (logger LogrusLogger) WithField(key string, value interface{}) Interface {
	return logrusEntry{
		Entry: logger.Logger.WithField(key, value),
	}
}

// WithFields append fields to log entry
func (logger *LogrusLogger) WithFields(fields Fields) Interface {
	return logrusEntry{
		Entry: logger.Logger.WithFields(map[string]interface{}(fields)),
	}
}

// LogrusEntry logrus log entry
type logrusEntry struct {
	*logrus.Entry
}

// WithField log with field
func (logger logrusEntry) WithField(key string, value interface{}) Interface {
	return logrusEntry{
		Entry: logger.Entry.WithField(key, value),
	}
}

// WithFields log with fields
func (logger logrusEntry) WithFields(fields Fields) Interface {
	return logrusEntry{
		Entry: logger.Entry.WithFields(map[string]interface{}(fields)),
	}
}
