package util

import (
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
	request, err := jiraClient.CreateRequest("GET", "/rest/api/latest/issue/WAL-1", nil)
	assert.NoError(err)
	assert.NotEmpty(request)
}
