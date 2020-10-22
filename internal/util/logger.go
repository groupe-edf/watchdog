package util

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/sirupsen/logrus"
)

// GetLogger create logger instance
func GetLogger(options *config.Options) *logrus.Logger {
	logger := logrus.New()
	logLevel, _ := logrus.ParseLevel(options.LogsLevel)
	logger.SetLevel(logLevel)
	logger.SetReportCaller(true)
	if options.LogsPath != "" {
		logFile, err := os.OpenFile(filepath.Clean(options.LogsPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			logger.Fatalf("Failed to open log file: %v", err)
		}
		logger.SetOutput(logFile)
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
	return logger
}

// PrintBanner Print watchdog banner
func PrintBanner(ctx context.Context, options *config.Options) error {
	t, err := template.New("watchdog").Parse(config.Banner)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Options":   options,
		"BuildInfo": version.GetBuildInfo(),
	}
	return t.Execute(os.Stdout, data)
}
