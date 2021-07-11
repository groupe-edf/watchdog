package server

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type BackupOptions struct {
	BackupDirectory string
}

func Backup(options BackupOptions) error {
	path := filepath.Join(options.BackupDirectory, "backup", time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(path, os.FileMode(0744)); err != nil {
		return errors.New("unable to create backup file")
	}
	return nil
}
