package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (api *API) GetLeaks(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindLeaks(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "leaks").
		SetHeader("Content-Range", fmt.Sprintf("leaks %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

func (api *API) GetLeak(r *http.Request) response.Response {
	vars := mux.Vars(r)
	leakID, err := strconv.ParseInt(vars["leak_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	leak, err := api.store.FindLeakByID(leakID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, leak)
}
