package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (api *API) DeleteAnalysis(r *http.Request) response.Response {
	vars := mux.Vars(r)
	analusisID, err := uuid.Parse(vars["analysis_id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.store.DeleteAnalysis(analusisID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "ANALYSIS SUCCESSFULLY DELETED",
	})
}

func (api *API) GetAnalysis(r *http.Request) response.Response {
	vars := mux.Vars(r)
	analusisID, err := uuid.Parse(vars["analysis_id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	analysis, err := api.store.FindAnalysisByID(&analusisID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, analysis)
}

func (api *API) GetAnalyzes(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindAnalyzes(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "analyzes").
		SetHeader("Content-Range", fmt.Sprintf("analyzes %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}
