package git

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

type AffectedFile struct {
	SrcMode string
	DstMode string
	SrcSHA  string
	DstSHA  string
	Status  string
	SrcPath string
	DstPath string
}

type AffectedFilesOptions struct {
	DiffFilter  string
	Env         []string
	NewCommitID string
	OldCommitID string
}

func GetAffectedFiles(repository *models.Repository, options *AffectedFilesOptions) (entries []*AffectedFile, err error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
		_ = writer.Close()
	}()
	entries = make([]*AffectedFile, 0, 10)
	args := []string{"diff-tree", "--raw", "--no-commit-id", "-r", "-z"}
	if options.DiffFilter != "" {
		args = append(args, fmt.Sprintf("--diff-filter=%s", options.DiffFilter))
	}
	args = append(args, options.OldCommitID, options.NewCommitID)
	err = NewCommand(context.Background(), args...).
		Run(&RunOptions{
			Env:    options.Env,
			Dir:    repository.Storage,
			Stdout: writer,
			PipelineFunc: func(_ context.Context, _ context.CancelFunc) error {
				_ = writer.Close()
				defer func() {
					_ = reader.Close()
				}()
				parser := NewParser(reader)
				for {
					diff, err := parser.NextDiff()
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("read diff: %v", err)
					}
					entries = append(entries, diff)
				}
				return nil
			},
		})
	return entries, err
}
