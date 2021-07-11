package models

import "github.com/google/uuid"

type WorkspaceMember struct {
	WorkspaceID *uuid.UUID `json:"workspace_id"`
	UserID      *uuid.UUID `json:"user_id"`
}
