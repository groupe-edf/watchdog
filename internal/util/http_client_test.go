package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHTTPClient(t *testing.T) {
	assert := assert.New(t)
	httpClient, err := CreateHTTPClient(nil, "http://localhost:8080")
	assert.NoError(err)
	assert.NotEmpty(httpClient)
}

func TestCreateRequest(t *testing.T) {
	assert := assert.New(t)
	jiraClient, err := CreateHTTPClient(nil, "http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	var endpoint = fmt.Sprintf("/rest/api/latest/issue/WAL-1")
	request, err := jiraClient.CreateRequest("GET", endpoint, nil)
	assert.NoError(err)
	assert.NotEmpty(request)
}
