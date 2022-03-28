package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/services"
	"github.com/groupe-edf/watchdog/pkg/query"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	var requestPayload = []byte(`{
		"email":"habib.maalem@gmail.com",
		"first_name": "Habib",
		"last_name": "MAALEM",
		"password": "watchdog"
	}`)
	t.Run("it returns 200 on /api/v1/authentication/register", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/api/v1/authentication/register", bytes.NewBuffer(requestPayload))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Cleanup(func() {
		api.store.DeleteUsers(&query.Query{})
	})
}

func TestRegisterExistingUser(t *testing.T) {
	assert := assert.New(t)
	userService := services.NewUserService(api.store)
	userService.CreateUser(&models.User{
		Email:    "habib.maalem@gmail.com",
		Password: "watchdog",
	})
	var requestPayload = []byte(`{
		"email":"habib.maalem@gmail.com",
		"first_name": "Habib",
		"last_name": "MAALEM",
		"password": "watchdog"
	}`)
	response := performRequest(router, http.MethodPost, "/api/v1/authentication/register", requestPayload)
	assert.Equal(http.StatusConflict, response.Code)
	assert.JSONEq(`{"detail":"USER_ALREADY_EXISTS","status":409}`, response.Body.String())
	t.Cleanup(func() {
		api.store.DeleteUsers(&query.Query{})
	})
}

func TestAuthentication(t *testing.T) {
	userService := services.NewUserService(api.store)
	userService.CreateUser(&models.User{
		Email:    "habib.maalem@gmail.com",
		Password: "watchdog",
	})
	var requestPayload = []byte(`{
		"email":"habib.maalem@gmail.com",
		"password": "watchdog"
	}`)
	response := performRequest(router, http.MethodPost, "/api/v1/authentication/login", requestPayload)
	assert.Equal(t, http.StatusOK, response.Code)
	var responsePayload map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &responsePayload)
	if responsePayload["token"] == "" {
		t.Errorf("Expected authentication token. Got %s", responsePayload["token"])
	}
	t.Cleanup(func() {
		api.store.DeleteUsers(&query.Query{})
	})
}

func TestBasicAuthentication(t *testing.T) {

}
