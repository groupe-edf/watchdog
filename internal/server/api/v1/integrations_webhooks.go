package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/util"
	"github.com/xanzy/go-gitlab"
)

type Webhook struct {
	Secret         string
	EventsToAccept []gitlab.EventType
}

type InstallWebhookCommand struct {
	IntegrationID int64 `json:"integration_id"`
	GroupID       int   `json:"group_id"`
}

func (api *API) InstallWebhook(r *http.Request) response.Response {
	var command InstallWebhookCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	integration, err := api.store.FindIntegrationByID(command.IntegrationID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	apiToken, err := util.Decrypt(integration.APIToken, api.options.Server.Security.MasterKey)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	client, err := gitlab.NewClient(apiToken, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", integration.InstanceURL)))
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	token := models.GenerateToken(32)
	originalWebhook, _, err := client.Groups.AddGroupHook(
		command.GroupID,
		&gitlab.AddGroupHookOptions{
			Token: &token,
			URL:   gitlab.String("https://watchdog.carecrute.com/api/v1/webhooks"),
		},
	)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	webhook, err := api.store.AddWebhook(&models.Webhook{
		IntegrationID: integration.ID,
		GroupID:       originalWebhook.GroupID,
		Token:         token,
		URL:           originalWebhook.URL,
		WebhookID:     originalWebhook.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, webhook)
}

func (api *API) HandleWebhook(r *http.Request) response.Response {
	return response.JSON(http.StatusOK, map[string]string{
		"message": "webhook successfully installer",
	})
}
