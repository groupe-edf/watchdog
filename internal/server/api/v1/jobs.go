package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (api *API) CancelJob(r *http.Request) response.Response {
	vars := mux.Vars(r)
	_, err := strconv.ParseInt(vars["job_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.Empty(http.StatusNoContent)
}

func (api *API) GetQueues(r *http.Request) response.Response {
	return response.Empty(http.StatusNoContent)
}

func (api *API) GetJobs(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindJobs(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "jobs").
		SetHeader("Content-Range", fmt.Sprintf("jobs %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}
