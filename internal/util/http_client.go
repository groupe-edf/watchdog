package util

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient data structure
type HTTPClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// CreateHTTPClient create new http client
func CreateHTTPClient(httpClient *http.Client, baseURL string) (*HTTPClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	client := &HTTPClient{
		baseURL:    parsedBaseURL,
		httpClient: httpClient,
	}
	return client, nil
}

// CreateRequest create raw request
func (client *HTTPClient) CreateRequest(method string, endpointURL string, body io.Reader) (*http.Request, error) {
	relativeURL, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	relativeURL.Path = strings.TrimLeft(relativeURL.Path, "/")
	endpoint := client.baseURL.ResolveReference(relativeURL)
	req, err := http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Do excute http request
func (client *HTTPClient) Do(req *http.Request, mapper interface{}) (*http.Response, error) {
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
