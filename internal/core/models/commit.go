package models

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

const GitTimeLayout = "Mon Jan _2 15:04:05 2006 -0700"

type Iterator[Model any] interface {
	Close()
	ForEach(func(*Model) error) error
	Next() (*Model, error)
}

type Commit struct {
	Author     *Signature          `json:"author"`
	Body       string              `json:"body"`
	Committer  *Signature          `json:"committer"`
	Hash       string              `json:"hash"`
	Parents    []string            `json:"parents"`
	Repository *Repository         `json:"-"`
	Signature  *CommitGPGSignature `json:"-"`
	Subject    string              `json:"subject"`
	Tree       interface{}         `json:"-"`
}

type Signature struct {
	Date     time.Time `json:"date"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Timezone string    `json:"timezone"`
	When     time.Time `json:"when"`
}

type CommitGPGSignature struct {
	Signature string
	Payload   string
}

func (signature *Signature) Decode(b []byte) {
	commitSignature, _ := parseCommitSignature(b)
	signature.Email = commitSignature.Email
	signature.Name = commitSignature.Name
	signature.When = commitSignature.When
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%s <%s>", signature.Name, signature.Email)
}

func parseCommitSignature(line []byte) (signature *Signature, err error) {
	signature = &Signature{}
	emailStart := bytes.LastIndexByte(line, '<')
	emailEnd := bytes.LastIndexByte(line, '>')
	if emailStart == -1 || emailEnd == -1 || emailEnd < emailStart {
		return
	}
	signature.Name = string(line[:emailStart-1])
	signature.Email = string(line[emailStart+1 : emailEnd])
	hasTime := emailEnd+2 < len(line)
	if !hasTime {
		return
	}
	firstChar := line[emailEnd+2]
	if firstChar >= 48 && firstChar <= 57 {
		idx := bytes.IndexByte(line[emailEnd+2:], ' ')
		if idx < 0 {
			return
		}
		timestring := string(line[emailEnd+2 : emailEnd+2+idx])
		seconds, _ := strconv.ParseInt(timestring, 10, 64)
		signature.When = time.Unix(seconds, 0)
		idx += emailEnd + 3
		if idx >= len(line) || idx+5 > len(line) {
			return
		}
		timezone := string(line[idx : idx+5])
		timezoneHours, _ := strconv.ParseInt(timezone[0:3], 10, 64)
		timezoneMinutes, _ := strconv.ParseInt(timezone[3:], 10, 64)
		if timezoneHours < 0 {
			timezoneMinutes *= -1
		}
		tz := time.FixedZone("", int(timezoneHours*60*60+timezoneMinutes*60))
		signature.When = signature.When.In(tz)
	} else {
		signature.When, err = time.Parse(GitTimeLayout, string(line[emailEnd+2:]))
		if err != nil {
			return
		}
	}
	return signature, err
}
