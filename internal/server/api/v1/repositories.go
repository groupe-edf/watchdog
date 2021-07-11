package v1

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/backend"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/issue"
	"github.com/groupe-edf/watchdog/internal/security"
	"github.com/groupe-edf/watchdog/internal/server/store"
)

func (api *API) Analyze(r *http.Request) Response {
	var command *AnalyzeRepositoryCommand
	var issues []issue.Issue
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return Response{Error: err}
	}
	repositoryURL, err := url.Parse(command.RepositoryURL)
	if err != nil {
		return Response{Error: err}
	}
	// Save repisotory if not exists
	ID := uuid.New()
	repository, err := api.store.SaveRepository(&store.Repository{
		ID:            &ID,
		RepositoryURL: repositoryURL.String(),
	})
	if err != nil {
		return Response{Error: err}
	}
	client := backend.New()
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	api.logger.Infof("fetching repository `%s`", repositoryURL)
	err = client.Clone(ctx, command.RepositoryURL)
	if err != nil {
		return Response{Error: err}
	}
	api.Broadcast("repository_fetched", repository)
	logOptions := &git.LogOptions{}
	from := plumbing.NewHash(command.From)
	if from != plumbing.ZeroHash {
		logOptions.From = from
	}
	since, err := time.Parse(command.Since, "2006-01-02")
	if err == nil {
		logOptions.Since = &since
	}
	until, err := time.Parse(command.Since, "2006-01-02")
	if err == nil {
		logOptions.Until = &until
	}
	commitIter, err := client.Commits(ctx, logOptions)
	if err != nil {
		return Response{Error: err}
	}
	defer commitIter.Close()
	// Start analysis
	createdAt := time.Now()
	analysisID := uuid.New()
	analysis, err := api.store.SaveAnalysis(&store.Analysis{
		ID:           &analysisID,
		CreatedAt:    &createdAt,
		RepositoryID: repository.ID,
		StartedAt:    &createdAt,
		State:        "STARTED",
	})
	if err != nil {
		return Response{Error: err}
	}
	api.Broadcast("repository_saved", repository)
	api.Broadcast("analysis_started", analysis)
	securityScanner := security.NewRegexScanner(api.logger, &config.Security{})
	for {
		commit, err := commitIter.Next()
		if err == io.EOF {
			break
		}
		leaks := make([]security.Leak, 0)
		commitLeaks, _ := securityScanner.Scan(commit)
		if len(commitLeaks) > 0 {
			leaks = append(leaks, commitLeaks...)
		}
		issues = append(issues, issue.Issue{
			Author: commit.Author.String(),
			Commit: commit.Hash.String(),
			Email:  commit.Author.Email,
			Leaks:  leaks,
		})
	}
	go func(analysis *store.Analysis) {
		severity := issue.SeverityLow
		for _, data := range issues {
			ID := uuid.New()
			if data.Severity > severity {
				severity = data.Severity
			}
			err := api.store.SaveIssue(&ID, analysis.ID, data)
			if err != nil {
				api.logger.Fatal(err)
			}
		}
		finishedAt := time.Now()
		elapsed := time.Since(*analysis.StartedAt)
		analysis.Duration = elapsed
		analysis.FinishedAt = &finishedAt
		analysis.Severity = severity
		analysis.State = "FINISHED"
		analysis.TotalIssues = len(issues)
		api.store.SaveAnalysis(analysis)
		api.Broadcast("analysis_finished", analysis)
	}(analysis)
	return Response{}
}

func (api *API) GetRepositories(r *http.Request) Response {
	repositories, err := api.store.FindRepositories()
	if err != nil {
		return Response{Error: err}
	}
	return Response{
		Data: repositories,
	}
}

func (api *API) GetRepository(r *http.Request) Response {
	vars := mux.Vars(r)
	repositoryID, err := uuid.Parse(vars["repository_id"])
	if err != nil {
		return Response{Error: err}
	}
	repository, err := api.store.FindRepositoryByID(&repositoryID)
	if err != nil {
		return Response{Error: err}
	}
	return Response{
		Data: repository,
	}
}
