package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) GetProfile(r *http.Request) response.Response {
	user, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	profile, err := api.store.FindUserById(user.ID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, profile)
}

func (api *API) UpdateProfile(r *http.Request) response.Response {
	return response.Empty(http.StatusOK)
}

// ChangePasswordCommand swagger:parameters
//
// Change password command
type ChangePasswordCommand struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	// User password
	//
	// required: true
	// in: body
	Password string `json:"password" validate:"required"`
	// User password confirmation
	//
	// required: true
	// in: body
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

func (api *API) ChangePassword(r *http.Request) response.Response {
	currentUser, err := token.GetUser(r.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	var command *ChangePasswordCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	validate := validator.New()
	err = validate.Struct(command)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	defer r.Body.Close()
	user, err := api.store.FindUserById(currentUser.ID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.CurrentPassword)); err != nil {
		return response.Error(http.StatusInternalServerError, "", errors.New("incorrect current password"))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.MinCost)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	user.Password = string(hashedPassword)
	_, err = api.store.UpdatePassword(user)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "password successfully changed",
	})
}
