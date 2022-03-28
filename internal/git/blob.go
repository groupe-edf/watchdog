package git

import (
	"context"
)

type Blob struct {
	ID         string
	name       string
	repository *Repository
	size       int64
}

func (blob *Blob) Name() string {
	return blob.name
}

func (blob *Blob) Size(ctx context.Context) int64 {
	writer, reader, cancel := CatFileBatchCheck(ctx, blob.repository.Path())
	defer cancel()
	_, err := writer.Write([]byte(blob.ID + "\n"))
	if err != nil {
		return 0
	}
	_, _, blob.size, err = ReadBatchLine(reader)
	if err != nil {
		return 0
	}
	return blob.size
}
