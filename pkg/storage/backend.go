package storage

import (
	"errors"
	"io"
	"os"

	"github.com/groupe-edf/watchdog/internal/config"
)

var (
	// ErrUnsupportedBackend Unsupported backend storage error
	ErrUnsupportedBackend = errors.New("unsupported backend storage")
)

const (
	// LocalStorage local storage backend name
	LocalStorage string = "local"
)

// Backend storage interface
type Backend interface {
	FileExists(path string) (bool, error)
	Reader(path string) (*os.File, error)
	ReadFile(path string) ([]byte, error)
	RemoveFile(path string) error
	WriteFile(reader io.Reader, path string) (int64, error)
}

// New create new backend storage
func New(options config.Storage) (Backend, error) {
	switch options.Driver {
	case LocalStorage:
		return &Local{
			directory: options.Directory,
		}, nil
	}
	return nil, ErrUnsupportedBackend
}
