package models

import "time"

type Integration struct {
	ID           int64      `json:"id,omitempty"`
	APIToken     string     `json:"-"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	CreatedBy    *User      `json:"created_by,omitempty"`
	InstanceName string     `json:"instance_name,omitempty"`
	InstanceURL  string     `json:"instance_url,omitempty"`
	SyncedAt     *time.Time `json:"synced_at,omitempty"`
	SyncingError string     `json:"syncing_error,omitempty"`
}
