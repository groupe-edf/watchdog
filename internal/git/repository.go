package git

import "path/filepath"

type Repository struct {
	path string
}

// Path returns the path to `repo`.
func (repo *Repository) Path() string {
	return repo.path
}

func NewRepository(path string) (*Repository, error) {
	path = filepath.Clean(path)
	repository := &Repository{
		path: path,
	}
	return repository, nil
}
