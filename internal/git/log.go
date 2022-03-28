package git

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"time"

	"github.com/gitleaks/go-gitdiff/gitdiff"
)

func GetLog(ctx context.Context, options []string, repository *Repository) (<-chan *gitdiff.File, error) {
	args := []string{"-C", repository.Path(), "log", "-p", "-U0"}
	if len(options) > 0 {
		args = append(args, options...)
	} else {
		args = append(args, "--full-history", "--all")
	}
	cmd := exec.Command("git", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go handleErrors(stderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	time.Sleep(50 * time.Millisecond)
	return gitdiff.Parse(cmd, stdout)
}

func handleErrors(stderr io.ReadCloser) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {

	}
}
