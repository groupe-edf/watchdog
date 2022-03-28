package jira

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/groupe-edf/watchdog/internal/util"
)

// Jira an http client to consume Jira restful API
type Jira struct {
	Authentication *Authentication
	client         *util.HTTPClient
}

// BasicAuthTransport data structure
type BasicAuthTransport struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

// QueryOptions data structure
type QueryOptions struct {
	Fields        string `url:"fields,omitempty"`
	Expand        string `url:"expand,omitempty"`
	Properties    string `url:"properties,omitempty"`
	FieldsByKeys  bool   `url:"fieldsByKeys,omitempty"`
	UpdateHistory bool   `url:"updateHistory,omitempty"`
	ProjectKeys   string `url:"projectKeys,omitempty"`
}

// Issue data structure
type Issue struct {
	ID string `json:"id,omitempty" structs:"id,omitempty"`
}

// New create a Jira client
func New(serverURL string) (*Jira, error) {
	httpClient, err := util.CreateHTTPClient(nil, serverURL)
	if err != nil {
		return nil, err
	}
	jira := &Jira{
		Authentication: &Authentication{client: httpClient},
		client:         httpClient,
	}
	return jira, nil
}

// GetIssue get Jira issue details
func (jira *Jira) GetIssue(issueID string) (*Issue, error) {
	issue := &Issue{}
	options := &QueryOptions{Properties: "ID"}
	urlQuery, _ := query.Values(options)
	var endpoint = fmt.Sprintf("/rest/api/latest/issue/%s", issueID)
	req, err := jira.client.CreateRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = urlQuery.Encode()
	res, err := jira.client.Do(req, issue)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(issue)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
