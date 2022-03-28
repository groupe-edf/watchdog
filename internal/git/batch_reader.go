package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"runtime"
	"strconv"
	"strings"
)

type WriteCloserError interface {
	io.WriteCloser
	CloseWithError(err error) error
}

func CatFileBatchCheck(ctx context.Context, repositoryPath string) (WriteCloserError, *bufio.Reader, func()) {
	batchStdinReader, batchStdinWriter := io.Pipe()
	batchStdoutReader, batchStdoutWriter := io.Pipe()
	ctx, ctxCancel := context.WithCancel(ctx)
	closed := make(chan struct{})
	cancel := func() {
		ctxCancel()
		_ = batchStdoutReader.Close()
		_ = batchStdinWriter.Close()
		<-closed
	}
	go func() {
		<-ctx.Done()
		cancel()
	}()
	_, filename, _, _ := runtime.Caller(2)
	filename = strings.TrimPrefix(filename, callerPrefix)
	go func() {
		stderr := strings.Builder{}
		err := NewCommand(ctx, "cat-file", "--batch-check").
			Run(&RunOptions{
				Dir:    repositoryPath,
				Stdin:  batchStdinReader,
				Stdout: batchStdoutWriter,
				Stderr: &stderr,
			})
		if err != nil {
			_ = batchStdoutWriter.CloseWithError(ConcatenateError(err, (&stderr).String()))
			_ = batchStdinReader.CloseWithError(ConcatenateError(err, (&stderr).String()))
		} else {
			_ = batchStdoutWriter.Close()
			_ = batchStdinReader.Close()
		}
		close(closed)
	}()
	batchReader := bufio.NewReader(batchStdoutReader)
	return batchStdinWriter, batchReader, cancel
}

func ReadTreeID(reader *bufio.Reader, size int64) (string, error) {
	var id string
	var n int64
headerLoop:
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		n += int64(len(line))
		idx := bytes.Index(line, []byte{' '})
		if idx < 0 {
			continue
		}
		if string(line[:idx]) == "tree" {
			id = string(line[idx+1 : len(line)-1])
			break headerLoop
		}
	}
	discard := size - n + 1
	for discard > math.MaxInt32 {
		_, err := reader.Discard(math.MaxInt32)
		if err != nil {
			return id, err
		}
		discard -= math.MaxInt32
	}
	_, err := reader.Discard(int(discard))
	return id, err
}

func ParseTreeLine(reader *bufio.Reader, modeBytes, fileNameBytes, oidBytes []byte, _ int64) (mode int64, fileName []byte, oid string, n int, err error) {
	var readBytes []byte
	readBytes, err = reader.ReadSlice('\x00')
	if err != nil {
		return
	}
	idx := bytes.IndexByte(readBytes, ' ')
	if idx < 0 {
		return
	}
	n += idx + 1
	copy(modeBytes, readBytes[:idx])
	if len(modeBytes) >= idx {
		modeBytes = modeBytes[:idx]
	} else {
		modeBytes = append(modeBytes, readBytes[len(modeBytes):idx]...)
	}
	mode, err = strconv.ParseInt(string(modeBytes), 10, 32)
	if err != nil {
		return
	}
	readBytes = readBytes[idx+1:]
	copy(fileNameBytes, readBytes)
	if len(fileNameBytes) > len(readBytes) {
		fileNameBytes = fileNameBytes[:len(readBytes)]
	} else {
		fileNameBytes = append(fileNameBytes, readBytes[len(fileNameBytes):]...)
	}
	for err == bufio.ErrBufferFull {
		readBytes, err = reader.ReadSlice('\x00')
		fileNameBytes = append(fileNameBytes, readBytes...)
	}
	n += len(fileNameBytes)
	if err != nil {
		return
	}
	fileNameBytes = fileNameBytes[:len(fileNameBytes)-1]
	fileName = fileNameBytes
	idx = 0
	for idx < 20 {
		var read int
		read, err = reader.Read(oidBytes[idx:20])
		n += read
		if err != nil {
			return
		}
		idx += read
	}
	oid = fmt.Sprintf("%02x", oidBytes[:20])
	return mode, fileName, oid, n, err
}

func catBatchParseTreeEntries(parentTree *Tree, reader *bufio.Reader, size int64) ([]*TreeEntry, error) {
	fileNameBytes := make([]byte, 4096)
	modeBytes := make([]byte, 40)
	oidBytes := make([]byte, 40)
	entries := make([]*TreeEntry, 0, 10)
loop:
	for size > 0 {
		mode, fileName, oid, count, err := ParseTreeLine(reader, modeBytes, fileNameBytes, oidBytes, size)
		if err != nil {
			if err == io.EOF {
				break loop
			}
			return nil, err
		}
		size -= int64(count)
		entry := new(TreeEntry)
		entry.parentTree = parentTree
		switch mode {
		case 100644:
			entry.entryMode = EntryModeBlob
		case 100755:
			entry.entryMode = EntryModeExec
		case 120000:
			entry.entryMode = EntryModeSymlink
		case 160000:
			entry.entryMode = EntryModeCommit
		case 40000:
			entry.entryMode = EntryModeTree
		default:
			return nil, fmt.Errorf("unknown mode: %v", mode)
		}
		entry.ID = oid
		entry.name = string(fileName)
		entries = append(entries, entry)
	}
	if _, err := reader.Discard(1); err != nil {
		return entries, err
	}
	return entries, nil
}

var callerPrefix string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	callerPrefix = strings.TrimSuffix(filename, "internal/git/batch_reader.go")
}
