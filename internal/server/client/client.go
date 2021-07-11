package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/groupe-edf/watchdog/internal/models"
	apiResponse "github.com/groupe-edf/watchdog/internal/server/api/response"
)

// Option is a functional option for configuring the API client
type Option func(*Client) error

// APIKey allows overriding of API client apiKey for testing
func WithAPIKey(apiKey string) Option {
	return func(client *Client) error {
		client.apiKey = apiKey
		return nil
	}
}

// Token allows overriding of API client token for testing
func WithToken(token string) Option {
	return func(client *Client) error {
		client.token = token
		return nil
	}
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	token      string
}

func (client *Client) GetPolicies(ctx context.Context) ([]models.Policy, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf(`%s/policies?conditions="policies.enabled",eq,true`, client.baseURL), nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response := make([]models.Policy, 0)
	if err := client.sendRequest(request, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (client *Client) GetRules(ctx context.Context) ([]models.Rule, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf(`%s/rules?conditions="rules.enabled",eq,true`, client.baseURL), nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	response := make([]models.Rule, 0)
	if err := client.sendRequest(request, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (client *Client) parseOptions(options ...Option) error {
	for _, option := range options {
		err := option(client)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) sendRequest(request *http.Request, data interface{}) error {
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")
	if client.token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	}
	if client.apiKey != "" {
		request.Header.Set(models.AccessTokenHeader, client.apiKey)
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		var problem apiResponse.DefaultProblem
		if err = json.NewDecoder(response.Body).Decode(&problem); err != nil {
			return err
		}
		return errors.New(problem.ProblemTitle())
	}
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return err
	}
	return nil
}

func NewClient(baseURL string, options ...Option) *Client {
	client := &Client{
		baseURL: baseURL + "/api/v1",
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
	// Loop through each option
	for _, option := range options {
		// Call the option giving the instantiated *Client as the argument
		option(client)
	}
	return client
}
