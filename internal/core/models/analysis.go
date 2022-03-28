package models

import (
	"time"

	"github.com/google/uuid"
)

type AnalysisTrigger string
type AnalysisState string

const (
	// Analysis trigger
	ManualTrigger    AnalysisTrigger = "manual"
	ScheduledTrigger AnalysisTrigger = "scheduled"
	// Analysis states
	CanceledState   AnalysisState = "canceled"
	PendingState    AnalysisState = "pending"
	InProgressState AnalysisState = "in_progress"
	SuccessState    AnalysisState = "success"
	FailedState     AnalysisState = "failed"
)

type Analysis struct {
	ID             *uuid.UUID      `json:"id"`
	CreatedAt      *time.Time      `json:"created_at,omitempty"`
	CreatedBy      *User           `json:"created_by,omitempty"`
	Duration       time.Duration   `json:"duration,omitempty"`
	FinishedAt     *time.Time      `json:"finished_at,omitempty"`
	LastCommitHash string          `json:"last_commit_hash,omitempty"`
	LastCommitDate *time.Time      `json:"last_commit_date,omitempty"`
	Repository     *Repository     `json:"repository,omitempty"`
	Severity       Score           `json:"severity"`
	StartedAt      *time.Time      `json:"started_at,omitempty"`
	State          AnalysisState   `json:"state,omitempty"`
	StateMessage   string          `json:"state_message,omitempty"`
	TotalIssues    int             `json:"total_issues"`
	Trigger        AnalysisTrigger `json:"trigger"`
}

func (analysis *Analysis) Done() bool {
	if analysis.State == SuccessState || analysis.State == FailedState || analysis.State == CanceledState {
		return true
	}
	return false
}

func (analysis *Analysis) Failed(err error) {
	analysis.State = FailedState
	analysis.StateMessage = err.Error()
}

func (analysis *Analysis) Prepare(trigger AnalysisTrigger) {
	analysis.State = PendingState
	analysis.Trigger = trigger
}

func (analysis *Analysis) Start() {
	startedAt := time.Now()
	analysis.StartedAt = &startedAt
	analysis.State = InProgressState
}

func (analysis *Analysis) Complete(severity Score, totalIssues int) {
	finishedAt := time.Now()
	elapsed := time.Since(*analysis.StartedAt)
	analysis.Duration = elapsed
	analysis.FinishedAt = &finishedAt
	analysis.Severity = severity
	analysis.State = SuccessState
	analysis.TotalIssues = totalIssues
}

func NewAnalysis(repository *Repository, createdBy *uuid.UUID) *Analysis {
	analysisID := uuid.New()
	createdAt := time.Now()
	return &Analysis{
		ID:        &analysisID,
		CreatedAt: &createdAt,
		CreatedBy: &User{
			ID: createdBy,
		},
		Repository: repository,
	}
}
