package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

func Clone(ctx context.Context, from, to string, options CloneOptions) error {
	combinedArgs := make([]string, len(globalCommandArgs))
	copy(combinedArgs, globalCommandArgs)
	return CloneWithArgs(ctx, from, to, combinedArgs, options)
}

func CloneWithArgs(ctx context.Context, from, to string, args []string, opts CloneOptions) (err error) {
	toDir := path.Dir(to)
	if err = os.MkdirAll(toDir, os.ModePerm); err != nil {
		return err
	}
	cmd := NewCommandContextNoGlobals(ctx, args...)
	envs := os.Environ()
	if opts.Authentication != nil {
		envs = append(envs,
			fmt.Sprintf("USERNAME=%s", opts.Authentication.Username),
			fmt.Sprintf("PASSWORD=%s", opts.Authentication.Password),
		)
		cmd.AddArguments("-c", `credential.helper="!f() { printf 'username=%s\npassword=%s\n' "$USERNAME" "$PASSWORD" };f"`)
	}
	cmd.AddArguments("clone")
	if opts.Bare {
		cmd.AddArguments("--bare")
	}
	if opts.Depth > 0 {
		cmd.AddArguments("--depth", strconv.Itoa(opts.Depth))
	}
	if opts.Mirror {
		cmd.AddArguments("--mirror")
	}
	if opts.Quiet {
		cmd.AddArguments("--quiet")
	}
	if len(opts.Branch) > 0 {
		cmd.AddArguments("--branch", opts.Branch)
	}
	cmd.AddArguments(from, to)
	if opts.Timeout <= 0 {
		opts.Timeout = -1
	}
	stderr := new(bytes.Buffer)
	if err = cmd.Run(&RunOptions{
		Timeout: opts.Timeout,
		Env:     envs,
		Stdout:  io.Discard,
		Stderr:  stderr,
	}); err != nil {
		return ConcatenateError(err, stderr.String())
	}
	return nil
}

func SetCredentials(ctx context.Context, credentials *BasicAuthentication) error {
	credentialHelper := fmt.Sprintf(`!f() { sleep 1; echo "username=%s"; echo "password=%s"; };f`, credentials.Username, credentials.Password)
	cmd := NewCommandContextNoGlobals(ctx).AddArguments("config")
	cmd.AddArguments("credential.helper", credentialHelper)
	stderr := new(bytes.Buffer)
	if err := cmd.Run(&RunOptions{
		Stdout: io.Discard,
		Stderr: stderr,
	}); err != nil {
		return ConcatenateError(err, stderr.String())
	}
	return nil
}
