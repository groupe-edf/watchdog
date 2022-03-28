package git

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/groupe-edf/watchdog/internal/core/models"
)

const maxUnixCommitDate = 1 << 53

var (
	ErrCommitNotFound = errors.New("commit not found")
)

func GetCommit(ctx context.Context, repository *Repository, commitID string) (*models.Commit, error) {
	writer, reader, cancel := CatFileBatch(ctx, repository.Path())
	defer cancel()
	_, _ = writer.Write([]byte(commitID + "\n"))
	return GetCommitFromBatchReader(repository, reader, commitID)
}

func CommitFromReader(repository *Repository, commitID string, reader io.Reader) (*models.Commit, error) {
	commit := &models.Commit{
		Hash:      commitID,
		Author:    &models.Signature{},
		Committer: &models.Signature{},
	}
	payloadSB := new(strings.Builder)
	messageSB := new(strings.Builder)
	signatureSB := new(strings.Builder)
	message := false
	pgpsig := false
	bufferedReader, ok := reader.(*bufio.Reader)
	if !ok {
		bufferedReader = bufio.NewReader(reader)
	}
readLoop:
	for {
		line, err := bufferedReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_, _ = messageSB.Write(line)
			}
			_, _ = payloadSB.Write(line)
			break readLoop
		}
		if pgpsig {
			if len(line) > 0 && line[0] == ' ' {
				_, _ = signatureSB.Write(line[1:])
				continue
			} else {
				pgpsig = false
			}
		}
		if !message {
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) == 0 {
				message = true
				_, _ = payloadSB.Write(line)
				continue
			}
			split := bytes.SplitN(trimmed, []byte{' '}, 2)
			var data []byte
			if len(split) > 1 {
				data = split[1]
			}
			switch string(split[0]) {
			case "tree":
				commit.Tree = NewTree(repository, string(data))
				_, _ = payloadSB.Write(line)
			case "parent":
				commit.Parents = append(commit.Parents, string(data))
				_, _ = payloadSB.Write(line)
			case "author":
				commit.Author = &models.Signature{}
				commit.Author.Decode(data)
				_, _ = payloadSB.Write(line)
			case "committer":
				commit.Committer = &models.Signature{}
				commit.Committer.Decode(data)
				_, _ = payloadSB.Write(line)
			case "gpgsig":
				_, _ = signatureSB.Write(data)
				_ = signatureSB.WriteByte('\n')
				pgpsig = true
			}
		} else {
			_, _ = messageSB.Write(line)
			_, _ = payloadSB.Write(line)
		}
	}
	commit.Subject = messageSB.String()
	commit.Signature = &models.CommitGPGSignature{
		Signature: signatureSB.String(),
		Payload:   payloadSB.String(),
	}
	if len(commit.Signature.Signature) == 0 {
		commit.Signature = nil
	}
	return commit, nil
}

func GetCommitFromBatchReader(repository *Repository, reader *bufio.Reader, commitID string) (*models.Commit, error) {
	_, entryType, size, err := ReadBatchLine(reader)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrCommitNotFound
		}
		return nil, err
	}
	switch entryType {
	case "commit":
		commit, err := CommitFromReader(repository, commitID, io.LimitReader(reader, size))
		if err != nil {
			return nil, err
		}
		_, err = reader.Discard(1)
		if err != nil {
			return nil, err
		}
		return commit, nil
	default:
		_, err = reader.Discard(int(size) + 1)
		if err != nil {
			return nil, err
		}
		return nil, ErrCommitNotFound
	}
}

func parseSignature(line string) *models.Signature {
	signature := &models.Signature{}
	splitName := strings.SplitN(line, "<", 2)
	signature.Name = strings.TrimSuffix(splitName[0], " ")
	if len(splitName) < 2 {
		return signature
	}
	line = splitName[1]
	splitEmail := strings.SplitN(line, ">", 2)
	if len(splitEmail) < 2 {
		return signature
	}
	signature.Email = splitEmail[0]
	secSplit := strings.Fields(splitEmail[1])
	if len(secSplit) < 1 {
		return signature
	}
	seconds, err := strconv.ParseInt(secSplit[0], 10, 64)
	if err != nil || seconds > maxUnixCommitDate || seconds < 0 {
		seconds = time.Now().Unix()
	}
	signature.Date = time.Unix(seconds, 0)
	if len(secSplit) == 2 {
		signature.Timezone = secSplit[1]
	}
	return signature
}

func subjectFromBody(body []byte) []byte {
	return bytes.TrimRight(bytes.SplitN(body, []byte("\n"), 2)[0], "\r\n")
}
