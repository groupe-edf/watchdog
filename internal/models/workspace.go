package models

import "github.com/google/uuid"

type Workspace struct {
	ID          *uuid.UUID `json:"id"`
	CreatedAt   int64      `json:"create_at"`
	DisplayName string     `json:"display_name"`
	Name        string     `json:"name"`
}
