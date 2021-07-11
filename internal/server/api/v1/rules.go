package v1

import (
	"net/http"

	"github.com/groupe-edf/watchdog/internal/security"
)

func (api *API) GetRules(r *http.Request) Response {
	return Response{
		Data: make([]security.Rule, 0),
	}
}
