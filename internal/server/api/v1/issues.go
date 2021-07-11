package v1

import "net/http"

func (api *API) GetIssues(r *http.Request) Response {
	issues, err := api.store.FindIssues()
	if err != nil {
		return Response{Error: err}
	}
	return Response{
		Data: issues,
	}
}
