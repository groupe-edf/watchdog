package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/services"
	"github.com/stretchr/testify/assert"
)

func TestFetchPoliciesWithBasicAuthentication(t *testing.T) {
	assert := assert.New(t)
	user := &models.User{
		Email:    "habib.maalem@gmail.com",
		Password: "watchdog",
	}
	userService := services.NewUserService(api.store)
	userService.CreateUser(user)
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/policies", nil)
	request.Header.Set("Authorization", "Basic "+basicAuth(user.Email, user.Password))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(http.StatusOK, response.Code)
}
