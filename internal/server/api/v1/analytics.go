package v1

import (
	"net/http"

	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/query"
)

func (api *API) GetAnalytics(r *http.Request) response.Response {
	err := api.store.RefreshAnalytics()
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	analytics, err := api.store.GetAnalytics()
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	leakCountBySeverity, err := api.store.GetLeakCountBySeverity()
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	totalItems := make(map[string]int)
	for _, containerType := range []string{"repositories", "repositories_analyzes", "repositories_issues", "repositories_leaks"} {
		count, err := api.store.Count(containerType, &query.Query{})
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		totalItems[containerType] = count
	}
	return response.JSON(http.StatusOK, map[string]interface{}{
		"total_items":            totalItems,
		"leak_count":             analytics,
		"leak_count_by_severity": leakCountBySeverity,
	})
}
