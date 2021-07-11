package models

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// Leak data struct
type Leak struct {
	ID          int64      `json:"id"`
	AnalysisID  uuid.UUID  `json:"-"`
	Author      string     `json:"author,omitempty"`
	CommitHash  string     `json:"commit_hash,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	File        string     `json:"file"`
	Line        string     `json:"line"`
	LineNumber  int        `json:"line_number"`
	Occurence   int        `json:"occurence"`
	Offender    string     `json:"offender"`
	Remediation string     `json:"remediation,omitempty"`
	Repository  Repository `json:"repository"`
	Rule        Rule       `json:"rule,omitempty"`
	SecretHash  string     `json:"secret_hash"`
	Severity    Severity   `json:"severity"`
	Tags        []string   `json:"tags,omitempty"`
}

func GenerateHash(args ...string) string {
	hasher := sha1.New()
	buffer := &bytes.Buffer{}
	gob.NewEncoder(buffer).Encode(args)
	hasher.Write(buffer.Bytes())
	return hex.EncodeToString(hasher.Sum(nil))
}
