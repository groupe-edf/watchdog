package models

import (
	"sync"
	"time"
)

type Job struct {
	ID         int64      `json:"id"`
	Args       []byte     `json:"args,omitempty"`
	ErrorCount int32      `json:"error_count"`
	LastError  string     `json:"last_error,omitempty"`
	Priority   int16      `json:"priority"`
	Queue      string     `json:"queue"`
	StartedAt  time.Time  `json:"started_at"`
	Type       string     `json:"type"`
	Deleted    bool       `json:"deleted"`
	Lock       sync.Mutex `json:"-"`
}
