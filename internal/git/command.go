package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	GitExecutable = "git"
	// globalCommandArgs global command args for external package setting
	globalCommandArgs []string
	// defaultCommandExecutionTimeout default command execution timeout duration
	defaultCommandExecutionTimeout = 360 * time.Second
)

// DefaultLocale is the default LC_ALL to run git commands in.
const DefaultLocale = "C"

type Command struct {
	args             []string
	executable       string
	globalArgsLength int
	parentContext    context.Context
}

type RunOptions struct {
	Env            []string
	Timeout        time.Duration
	Dir            string
	Stdout, Stderr io.Writer
	Stdin          io.Reader
	PipelineFunc   func(context.Context, context.CancelFunc) error
}

type RunStdError interface {
	error
	Unwrap() error
	Stderr() string
	IsExitCode(code int) bool
}

type runStdError struct {
	err    error
	stderr string
	errMsg string
}

func (r *runStdError) Error() string {
	if r.errMsg == "" {
		r.errMsg = ConcatenateError(r.err, r.stderr).Error()
	}
	return r.errMsg
}

func (r *runStdError) Unwrap() error {
	return r.err
}

func (r *runStdError) Stderr() string {
	return r.stderr
}

func (r *runStdError) IsExitCode(code int) bool {
	var exitError *exec.ExitError
	if errors.As(r.err, &exitError) {
		return exitError.ExitCode() == code
	}
	return false
}

// AddArguments adds new argument(s) to the command.
func (c *Command) AddArguments(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) RunStdString(opts *RunOptions) (stdout, stderr string, runErr RunStdError) {
	stdoutBytes, stderrBytes, err := c.RunStdBytes(opts)
	stdout = string(stdoutBytes)
	stderr = string(stderrBytes)
	if err != nil {
		return stdout, stderr, &runStdError{err: err, stderr: stderr}
	}
	// even if there is no err, there could still be some stderr output, so we just return stdout/stderr as they are
	return stdout, stderr, nil
}

func (c *Command) RunStdBytes(opts *RunOptions) (stdout, stderr []byte, runErr RunStdError) {
	if opts == nil {
		opts = &RunOptions{}
	}
	if opts.Stdout != nil || opts.Stderr != nil {
		panic("stdout and stderr field must be nil when using RunStdBytes")
	}
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	opts.Stdout = stdoutBuf
	opts.Stderr = stderrBuf
	err := c.Run(opts)
	stderr = stderrBuf.Bytes()
	if err != nil {
		return nil, stderr, &runStdError{err: err, stderr: string(stderr)}
	}
	return stdoutBuf.Bytes(), stderr, nil
}

func (c *Command) Run(opts *RunOptions) error {
	if opts == nil {
		opts = &RunOptions{}
	}
	if opts.Timeout <= 0 {
		opts.Timeout = defaultCommandExecutionTimeout
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, c.executable, c.args...)
	if opts.Env == nil {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = opts.Env
	}
	cmd.Env = append(cmd.Env, CommonGitCmdEnvs()...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	cmd.Stdin = opts.Stdin
	fmt.Println(cmd.String())
	if err := cmd.Start(); err != nil {
		return err
	}
	if opts.PipelineFunc != nil {
		err := opts.PipelineFunc(ctx, cancel)
		if err != nil {
			cancel()
			_ = cmd.Wait()
			return err
		}
	}
	if err := cmd.Wait(); err != nil && ctx.Err() != context.DeadlineExceeded {
		return err
	}
	return ctx.Err()
}

func (c *Command) String() string {
	if len(c.args) == 0 {
		return c.executable
	}
	return fmt.Sprintf("%s %s", c.executable, strings.Join(c.args, " "))
}

func NewCommand(ctx context.Context, args ...string) *Command {
	cargs := make([]string, len(globalCommandArgs))
	copy(cargs, globalCommandArgs)
	return &Command{
		executable:       GitExecutable,
		args:             append(cargs, args...),
		parentContext:    ctx,
		globalArgsLength: len(globalCommandArgs),
	}
}

func NewCommandContextNoGlobals(ctx context.Context, args ...string) *Command {
	return &Command{
		executable:    GitExecutable,
		args:          args,
		parentContext: ctx,
	}
}

func CommonGitCmdEnvs() []string {
	return []string{
		fmt.Sprintf("LC_ALL=%s", DefaultLocale),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_NO_REPLACE_OBJECTS=1",
	}
}

// ConcatenateError concatenats an error with stderr string
func ConcatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%w - %s", err, stderr)
}
