package models

import (
	"time"
)

type AnalysisResult struct {
	Commit      Commit
	ElapsedTime time.Duration
	Issues      []Issue
}
