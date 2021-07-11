package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
)

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func (api *API) OAuthLogin(r *http.Request) response.Response {
	var currentProvider config.OAuthProvider
	for _, provider := range api.options.Server.OAuthProviders {
		if provider.Name == "oauth_github" {
			currentProvider = *provider
			break
		}
	}
	code := r.FormValue("code")
	if code != "" {
		requestURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", currentProvider.ApplicationID, currentProvider.ApplicationSecret, code)
		req, err := http.NewRequest(http.MethodPost, requestURL, nil)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		req.Header.Set("Accept", "application/json")
		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		defer res.Body.Close()
		var token OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
			return response.Error(http.StatusInternalServerError, "", err)
		}
		return response.Empty(http.StatusOK).SetHeader("Location", "/welcome.html?access_token="+token.AccessToken)
	} else {
		return response.Error(http.StatusInternalServerError, "missing code parameter", nil)
	}
}
