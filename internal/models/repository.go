package models

import (
	"time"

	"github.com/google/uuid"
)

type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityPublic  Visibility = "public"
)

type Repository struct {
	ID               *uuid.UUID           `json:"id"`
	Analytics        *RepositoryAnalytics `json:"analytics,omitempty"`
	CreatedAt        *time.Time           `json:"created_at,omitempty"`
	CreatedBy        *uuid.UUID           `json:"created_by,omitempty"`
	EnableMonitoring bool                 `json:"enable_monitoring,omitempty"`
	Integration      Integration          `json:"integration,omitempty"`
	LastAnalysis     *Analysis            `json:"last_analysis,omitempty"`
	RepositoryURL    string               `json:"repository_url"`
	Visibility       Visibility           `json:"visibility,omitempty"`
}

type RepositoryAnalytics struct {
	AnalysisCount int64 `json:"analysis_count"`
	CommitCount   int64 `json:"commit_count"`
	IssueCount    int64 `json:"issue_count"`
	LeakCount     int64 `json:"leak_count"`
}
