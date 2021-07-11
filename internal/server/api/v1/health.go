package v1

import (
	"net/http"

	"github.com/groupe-edf/watchdog/internal/server/api/response"
)

func (api *API) Health(r *http.Request) response.Response {
	return response.Empty(http.StatusOK)
}
