package git

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func CatFileBatch(ctx context.Context, repositoryPath string) (io.WriteCloser, *bufio.Reader, func()) {
	batchStdinReader, batchStdinWriter := io.Pipe()
	batchStdoutReader, batchStdoutWriter := io.Pipe()
	ctx, ctxCancel := context.WithCancel(ctx)
	closed := make(chan struct{})
	cancel := func() {
		ctxCancel()
		_ = batchStdinWriter.Close()
		_ = batchStdoutReader.Close()
		<-closed
	}
	go func() {
		stderr := strings.Builder{}
		err := NewCommand(ctx, "cat-file", "--batch").Run(&RunOptions{
			Dir:    repositoryPath,
			Stdin:  batchStdinReader,
			Stdout: batchStdoutWriter,
			Stderr: &stderr,
		})
		if err != nil {
			_ = batchStdoutWriter.CloseWithError(err)
			_ = batchStdinReader.CloseWithError(err)
		} else {
			_ = batchStdoutWriter.Close()
			_ = batchStdinReader.Close()
		}
		close(closed)
	}()
	batchReader := bufio.NewReaderSize(batchStdoutReader, 32*1024)
	return batchStdinWriter, batchReader, cancel
}

func ReadBatchLine(reader *bufio.Reader) (sha []byte, typ string, size int64, err error) {
	typ, err = reader.ReadString('\n')
	if err != nil {
		return
	}
	if len(typ) == 1 {
		typ, err = reader.ReadString('\n')
		if err != nil {
			return
		}
	}
	idx := strings.IndexByte(typ, ' ')
	if idx < 0 {
		err = fmt.Errorf("commit not found %s", sha)
		return
	}
	sha = []byte(typ[:idx])
	typ = typ[idx+1:]
	idx = strings.IndexByte(typ, ' ')
	if idx < 0 {
		err = fmt.Errorf("commit not found %s", sha)
		return
	}
	sizeStr := typ[idx+1 : len(typ)-1]
	typ = typ[:idx]
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	return
}
