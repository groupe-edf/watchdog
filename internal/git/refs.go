package git

import (
	"context"
	"fmt"

	"github.com/groupe-edf/watchdog/internal/github"
)

type RefGroupSymbol string

type refSeen struct {
	github.Reference
	walked bool
	groups []RefGroupSymbol
}

func GetRefsFiltered(repositoryPath string) ([]*github.Reference, error) {
	repository, err := github.NewRepository(repositoryPath)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	refIter, err := repository.NewReferenceIter(ctx)
	errChan := make(chan error, 1)
	var _ []refSeen
	go func() {
		errChan <- func() error {
			for {
				ref, ok, err := refIter.Next()
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}
				fmt.Print(ref)
			}
		}()
	}()
	return nil, nil
}
