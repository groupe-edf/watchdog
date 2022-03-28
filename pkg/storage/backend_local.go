package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Local local storage backend
type Local struct {
	directory string
}

// FileExists check if file exists
func (backend *Local) FileExists(path string) (bool, error) {
	_, err := os.Stat(filepath.Join(backend.directory, path))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Reader retrun file reader
func (backend *Local) Reader(path string) (*os.File, error) {
	file, err := os.Open(filepath.Join(backend.directory, path))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ReadFile read file content from disk
func (backend *Local) ReadFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(filepath.Join(backend.directory, path))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// RemoveFile remove file
func (backend *Local) RemoveFile(path string) error {
	if err := os.Remove(filepath.Join(backend.directory, path)); err != nil {
		return err
	}
	return nil
}

// SetDirectory set base directory
func (backend *Local) SetDirectory(directory string) {
	backend.directory = directory
}

// WriteFile write file to disk
func (backend *Local) WriteFile(reader io.Reader, path string) (int64, error) {
	path = filepath.Join(backend.directory, path)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return 0, err
	}
	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	written, err := io.Copy(writer, reader)
	if err != nil {
		return written, err
	}
	return written, nil
}
