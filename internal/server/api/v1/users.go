package v1

import (
	"fmt"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (api *API) GetUsers(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindUsers(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "users").
		SetHeader("Content-Range", fmt.Sprintf("users %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}
