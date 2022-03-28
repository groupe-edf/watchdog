package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (api *API) GetAccessTokens(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindAccessTokens(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "access_tokens").
		SetHeader("Content-Range", fmt.Sprintf("access_tokens %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

type AddAccessToken struct {
	Name string `json:"name"`
}

func (api *API) SaveAccessToken(r *http.Request) response.Response {
	var command *AddAccessToken
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	user, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	token := models.NewAccessToken(command.Name, user.ID)
	token, err = api.store.SaveAccessToken(token)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, token)
}
