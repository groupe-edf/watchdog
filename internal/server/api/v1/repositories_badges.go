package v1

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
)

var (
	badgeColorMap = map[models.AnalysisState]string{
		models.SuccessState:    "#4c1",
		models.FailedState:     "#e05d44",
		models.InProgressState: "#dfb317",
	}
)
var badgeTemplate = strings.TrimSpace(`
<svg xmlns="http://www.w3.org/2000/svg" width="90" height="20">
	<title>build: {{ .State }}</title>
	<defs>
		<linearGradient id="workflow-fill" x1="50%" y1="0%" x2="50%" y2="100%">
			<stop stop-color="#444D56" offset="0%"/>
			<stop stop-color="#24292E" offset="100%"/>
		</linearGradient>
		<linearGradient id="state-fill" x1="50%" y1="0%" x2="50%" y2="100%">
			<stop stop-color="{{ .Color }}" offset="0%"/>
			<stop stop-color="{{ .Color }}" offset="100%"/>
		</linearGradient>
	</defs>
	<rect rx="3" width="90" height="20" fill="#555"/>
	<rect rx="3" x="40" width="50" height="20" fill="{{ .Color }}"/>
	<path fill="{{ .Color }}" d="M36 0h8v20h-8z"/>
	<rect rx="3" width="90" height="20" fill="url(#a)"/>
	<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
		<text x="18" y="15" fill="#010101" fill-opacity=".3">build</text>
		<text x="18" y="14">build</text>
		<text x="62" y="15" fill="#010101" fill-opacity=".3">{{ .State }}</text>
		<text x="62" y="14">{{ .State }}</text>
	</g>
</svg>
`)

func (api *API) GetRepositoryBadge(r *http.Request) response.Response {
	vars := mux.Vars(r)
	repositoryID, err := uuid.Parse(vars["repository_id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	repository, err := api.store.FindRepositoryByID(&repositoryID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	t := template.Must(template.New("watchdog").Parse(badgeTemplate))
	data := map[string]interface{}{
		"State": repository.LastAnalysis.State,
		"Color": badgeColorMap[repository.LastAnalysis.State],
	}
	var badge bytes.Buffer
	if err := t.Execute(&badge, data); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	badgeResponse := response.JSON(http.StatusOK, badge.Bytes())
	badgeResponse.SetHeader("Content-Type", "image/svg+xml")
	badgeResponse.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	return badgeResponse
}
