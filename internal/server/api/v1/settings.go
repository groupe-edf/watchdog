package v1

import (
	"net/http"

	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (api *API) GetSettings(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	settings, err := api.store.GetSettings(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	data := make(map[string]interface{})
	for _, setting := range settings {
		data[setting.SettingKey] = setting.CastValue()
	}
	return response.JSON(http.StatusOK, data)
}

func (api *API) SaveSettings(r *http.Request) response.Response {
	return response.Empty(http.StatusOK)
}
