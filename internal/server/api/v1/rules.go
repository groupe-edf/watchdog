package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/query"
)

type ToggleCommand struct {
	Enabled bool `json:"enabled"`
}

type EvaluatePatternCommand struct {
	Payload string `json:"payload"`
	Pattern string `json:"pattern"`
}

func (api *API) EvaluatePattern(r *http.Request) response.Response {
	var command *EvaluatePatternCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	_, err := regexp.Compile(command.Pattern)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "pattern successfully evaluated",
	})
}

func (api *API) GetRules(r *http.Request) response.Response {
	q := query.Parse(r.URL.Query())
	paginator, err := api.store.FindRules(q)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, paginator.Items).
		SetHeader("Accept-Ranges", "rules").
		SetHeader("Content-Range", fmt.Sprintf("rules %d-%d/%d", paginator.Query.Offset, paginator.Query.Limit, paginator.TotalItems))
}

func (api *API) ToggleRule(r *http.Request) response.Response {
	var command *ToggleCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	vars := mux.Vars(r)
	ruleID, err := strconv.ParseInt(vars["rule_id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	err = api.store.ToggleRule(ruleID, command.Enabled)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "rule successfully updated",
	})
}
