package v1

import (
	"fmt"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (api *API) GetCategories(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindCategories(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "categories").
		SetHeader("Content-Range", fmt.Sprintf("categories %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}
