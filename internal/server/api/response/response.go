package response

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response interface {
	Body() []byte
	Header() http.Header
	Status() int
}

type defaultResponse struct {
	body   *bytes.Buffer
	header http.Header
	status int
}

func (response *defaultResponse) Body() []byte {
	return response.body.Bytes()
}

func (response *defaultResponse) Header() http.Header {
	return response.header
}

func (response *defaultResponse) SetHeader(key, value string) *defaultResponse {
	response.header.Set(key, value)
	return response
}

func (response *defaultResponse) Status() int {
	return response.status
}

func Empty(status int) *defaultResponse {
	return Respond(status, nil)
}

func Error(status int, message string, err error) *defaultResponse {
	problem := NewStatusProblem(status)
	problem.Title = message
	problem.Detail = err.Error()
	response := JSON(status, problem)
	response.SetHeader("Content-Type", ProblemMediaType)
	return response
}

func JSON(status int, body interface{}) *defaultResponse {
	return Respond(status, body)
}

func Respond(status int, body interface{}) *defaultResponse {
	var responseBody []byte
	switch bodyType := body.(type) {
	case []byte:
		responseBody = bodyType
	case string:
		responseBody = []byte(bodyType)
	default:
		if body != nil {
			var err error
			responseBody, err = json.Marshal(body)
			if err != nil {
				return Error(500, "body json marshal", err)
			}
		}
	}
	return &defaultResponse{
		body:   bytes.NewBuffer(responseBody),
		header: make(http.Header),
		status: status,
	}
}

func Success() {}
