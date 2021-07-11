package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID          int64      `json:"id"`
	Content     string     `json:"content"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	NotableID   int64      `json:"notable_id"`
	NotableType string     `json:"notable_type"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
