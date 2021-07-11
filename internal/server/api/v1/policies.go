package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

func (api *API) GetPolicies(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindPolicies(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "policies").
		SetHeader("Content-Range", fmt.Sprintf("policies %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

func (api *API) GetPolicy(r *http.Request) response.Response {
	vars := mux.Vars(r)
	policyID, err := strconv.ParseInt(vars["policy_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	policy, err := api.store.FindPolicyByID(policyID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, policy)
}

func (api *API) TogglePolicy(r *http.Request) response.Response {
	var command *ToggleCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	vars := mux.Vars(r)
	policyID, err := strconv.ParseInt(vars["policy_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.store.TogglePolicy(policyID, command.Enabled)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "policy successfully updated",
	})
}
