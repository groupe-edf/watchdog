package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/authentication/token"
	"github.com/groupe-edf/watchdog/internal/server/job"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

type AnalyzeRepositoryCommand struct {
	CreatedBy        *uuid.UUID `json:"-"`
	EnableMonitoring bool       `json:"enable_monitoring"`
	From             string     `json:"from"`
	RepositoryURL    string     `json:"repository_url"`
	Since            string     `json:"since,omitempty"`
	Token            string     `json:"token,omitempty"`
	Until            string     `json:"until,omitempty"`
	Username         string     `json:"username,omitempty"`
}

func (api *API) Analyze(r *http.Request) response.Response {
	var command *AnalyzeRepositoryCommand
	var repository *models.Repository
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	user, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	createdAt := time.Now()
	command.CreatedBy = user.ID
	vars := mux.Vars(r)
	if vars["repository_id"] != "" {
		repositoryID, err := uuid.Parse(vars["repository_id"])
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		repository, err = api.store.FindRepositoryByID(&repositoryID)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
	} else {
		repositoryURL, err := url.Parse(command.RepositoryURL)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		// Save repisotory if not exists
		ID := uuid.New()
		repository, err = api.store.SaveRepository(&models.Repository{
			ID:               &ID,
			CreatedAt:        &createdAt,
			CreatedBy:        command.CreatedBy,
			EnableMonitoring: command.EnableMonitoring,
			RepositoryURL:    repositoryURL.String(),
			Visibility:       models.VisibilityPublic,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
	}
	if repository.LastAnalysis != nil && repository.LastAnalysis.State == models.InProgressState {
		return response.Error(http.StatusInternalServerError, "", errors.New("an analysis is already in execution"))
	}
	// Start analysis
	analysis := models.NewAnalysis(repository, command.CreatedBy)
	analysis.Prepare(models.ManualTrigger)
	analysis, err = api.store.SaveAnalysis(analysis)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	options := job.AnalyzeOptions{
		AnalysisID:    analysis.ID,
		CreatedAt:     &createdAt,
		CreatedBy:     command.CreatedBy,
		RepositoryID:  repository.ID,
		RepositoryURL: repository.RepositoryURL,
	}
	if repository.Integration.ID != 0 {
		options.IntegrationID = repository.Integration.ID
	}
	if repository.LastAnalysis != nil {
		options.From = repository.LastAnalysis.LastCommitHash
	}
	if command.From != "" {
		options.From = command.From
	}
	args, err := json.Marshal(options)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.workerPool.Enqueue(&models.Job{
		Args:  args,
		Queue: "default",
		Type:  "analyze_repository",
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, analysis)
}

func (api *API) DeleteRepository(r *http.Request) response.Response {
	vars := mux.Vars(r)
	repositoryID, err := uuid.Parse(vars["repository_id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.store.DeleteRepository(repositoryID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "repository successfully deleted",
	})
}

func (api *API) GetRepositories(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindRepositories(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "repositories").
		SetHeader("Content-Range", fmt.Sprintf("repositories %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

func (api *API) GetRepository(r *http.Request) response.Response {
	vars := mux.Vars(r)
	repositoryID, err := uuid.Parse(vars["repository_id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	repository, err := api.store.FindRepositoryByID(&repositoryID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, repository)
}
