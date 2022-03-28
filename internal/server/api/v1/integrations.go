package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/groupe-edf/watchdog/pkg/query"
	"github.com/xanzy/go-gitlab"
)

func (api *API) GetIntegrations(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindIntegrations(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "integrations").
		SetHeader("Content-Range", fmt.Sprintf("integrations %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

func (api *API) GetIntegration(r *http.Request) response.Response {
	vars := mux.Vars(r)
	integrationID, err := strconv.ParseInt(vars["integration_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	integration, err := api.store.FindIntegrationByID(integrationID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, integration)
}

type IntegrationGroup struct {
	ID         int    `json:"id"`
	Installed  bool   `json:"installed"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	WebhookURL string `json:"webhook_url"`
}

func (api *API) GetIntegrationGroup(r *http.Request) response.Response {
	vars := mux.Vars(r)
	integrationID, err := strconv.ParseInt(vars["integration_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	integration, err := api.store.FindIntegrationByID(integrationID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	token, err := util.Decrypt(integration.APIToken, api.options.Server.Security.MasterKey)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", integration.InstanceURL)))
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	options := &gitlab.ListGroupsOptions{}
	groups, _, err := client.Groups.ListGroups(options)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	q := &query.Query{}
	q.AddCondition(query.Condition{
		Field:    `"integrations_webhooks"."id"`,
		Operator: query.Equal,
		Value:    integration.ID,
	})
	webhooks, err := api.store.FindWebhooks(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	var integrationGroups []*IntegrationGroup
	for _, group := range groups {
		integrationGroup := &IntegrationGroup{
			ID:   group.ID,
			Name: group.Name,
			Path: group.Path,
		}
		for _, webhook := range webhooks.Items {
			if group.ID == webhook.GroupID {
				integrationGroup.Installed = true
				integrationGroup.WebhookURL = webhook.URL
			}
		}
		integrationGroups = append(integrationGroups, integrationGroup)
	}
	return response.JSON(http.StatusOK, integrationGroups)
}

func (api *API) DeleteIntegration(r *http.Request) response.Response {
	vars := mux.Vars(r)
	integrationID, err := strconv.ParseInt(vars["integration_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.store.DeleteIntegration(integrationID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "integration successfully deleted",
	})
}

type AddIntegration struct {
	APIToken     string `json:"api_token"`
	InstanceName string `json:"instance_name"`
	InstanceURL  string `json:"instance_url"`
}

func (api *API) SaveIntegration(r *http.Request) response.Response {
	var command *AddIntegration
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	user, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	createdAt := time.Now()
	token, err := util.Encrypt(command.APIToken, api.options.Server.Security.MasterKey)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	integration, err := api.store.SaveIntegration(&models.Integration{
		APIToken:  token,
		CreatedAt: &createdAt,
		CreatedBy: &models.User{
			ID: user.ID,
		},
		InstanceName: command.InstanceName,
		InstanceURL:  command.InstanceURL,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, integration)
}

func (api *API) SynchronizeIntegration(r *http.Request) response.Response {
	user, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	vars := mux.Vars(r)
	integrationID, err := strconv.ParseInt(vars["integration_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	integration, err := api.store.FindIntegrationByID(integrationID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	fmt.Print(integration.APIToken, api.options.Server.Security.MasterKey)
	token, err := util.Decrypt(integration.APIToken, api.options.Server.Security.MasterKey)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", integration.InstanceURL)))
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	options := &gitlab.ListProjectsOptions{
		Owned: gitlab.Bool(true),
	}
	currentDateTime := time.Now()
	projects, _, err := client.Projects.ListProjects(options)
	if err != nil {
		integration.SyncingError = err.Error()
		integration.SyncedAt = &currentDateTime
		integration, err = api.store.UpdateIntegration(integration)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		return response.JSON(http.StatusOK, integration)
	}
	for _, project := range projects {
		ID := uuid.New()
		var visibility models.Visibility
		if project.Visibility == gitlab.InternalVisibility {
			visibility = models.VisibilityPrivate
		} else {
			visibility = models.Visibility(project.Visibility)
		}
		_, err = api.store.SaveRepository(&models.Repository{
			ID:            &ID,
			CreatedAt:     &currentDateTime,
			CreatedBy:     user.ID,
			Integration:   integration,
			RepositoryURL: project.WebURL,
			Visibility:    visibility,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
	}
	integration.SyncedAt = &currentDateTime
	_, err = api.store.UpdateIntegration(integration)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("%d projects successfully imported", len(projects)),
	})
}
