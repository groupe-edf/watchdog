package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/query"
)

type ToggleCommand struct {
	Enabled bool `json:"enabled"`
}

type EvaluatePatternCommand struct {
	Payload         string `json:"payload"`
	Pattern         string `json:"pattern"`
	FindAllSubmatch bool   `json:"find_all_submatch"`
}

type MatchResultResponse struct {
	Matches    [][]string `json:"matches"`
	GroupsName []string   `json:"groups_name"`
}

func (api *API) EvaluatePattern(r *http.Request) response.Response {
	var matches [][]string
	var command *EvaluatePatternCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	pattern, err := regexp.Compile(command.Pattern)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	matchResult := &MatchResultResponse{}
	if command.FindAllSubmatch {
		matches = pattern.FindAllStringSubmatch(command.Payload, -1)
	} else {
		matches = [][]string{pattern.FindStringSubmatch(command.Payload)}
	}
	if len(matches) > 0 {
		matchResult.Matches = matches
		matchResult.GroupsName = pattern.SubexpNames()[1:]
	}
	return response.JSON(http.StatusOK, matchResult)
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

type AddRuleCommand struct {
	Description string          `json:"description"`
	DisplayName string          `json:"display_name"`
	Enabled     bool            `json:"enabled"`
	Pattern     string          `json:"pattern"`
	Severity    models.Severity `json:"severity"`
}

func (api *API) NewRule(r *http.Request) response.Response {
	var command *AddRuleCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	rule, err := api.store.SaveRule(&models.Rule{
		DisplayName: command.DisplayName,
		Description: command.Description,
		Enabled:     command.Enabled,
		Pattern:     command.Pattern,
		Severity:    command.Severity,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, rule)
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
