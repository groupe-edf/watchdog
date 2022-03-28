package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/services"
	"github.com/groupe-edf/watchdog/pkg/authentication"
	"github.com/groupe-edf/watchdog/pkg/authentication/provider"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
)

type LoginCommand struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterCommand struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type ForgotPasswordCommand struct {
	Email string `json:"email" validate:"required"`
}

func (api *API) Forgot(r *http.Request) response.Response {
	var command *ForgotPasswordCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	if err := validator.New().Struct(command); err != nil {
		return response.Error(http.StatusBadRequest, "", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "FORGOT_PASSWORD_EMAIL_SENT",
	})
}

func (api *API) Login(r *http.Request) response.Response {
	var command *LoginCommand
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	if err := validator.New().Struct(command); err != nil {
		return response.Error(http.StatusBadRequest, "", err)
	}
	authenticator := authentication.NewService(authentication.Options{})
	authenticator.AddProvider(provider.NewLDAPProvider(api.options.Server.LDAP))
	authenticator.AddProvider(provider.NewLocalProvider())
	identity, err := authenticator.Authenticate(command.Email, command.Password)
	if err != nil {
		return response.Error(http.StatusNotFound, "", models.ErrInvalidUsernameOrPassword)
	}
	if identity == nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	// Upsert user and set last login date
	id := uuid.New()
	currentDate := time.Now()
	user, err := api.store.SaveOrUpdateUser(&models.User{
		ID:        &id,
		CreatedAt: currentDate,
		Email:     identity.Email,
		FirstName: identity.FirstName,
		LastLogin: &currentDate,
		LastName:  identity.LastName,
		Provider:  string(identity.Provider),
		Username:  identity.Username,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	identity.ID = user.ID
	authenticationService := token.NewJWTService(token.JWTOptions{
		Secret: []byte("secret"),
	})
	token, err := authenticationService.Token(identity.ToClaims(token.Attributes{}))
	if err != nil {
		return response.Error(http.StatusInternalServerError, "invalid token", err)
	}
	return response.JSON(http.StatusOK, map[string]string{
		"email":      identity.Email,
		"first_name": identity.FirstName,
		"last_name":  identity.LastName,
		"token":      token,
	}).SetHeader(authenticationService.Options.JWTHeaderKey, token)
}

func (api *API) Register(r *http.Request) response.Response {
	var command *RegisterCommand
	var user *models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&command); err != nil {
		return response.Error(http.StatusInternalServerError, "", err)
	}
	defer r.Body.Close()
	user = &models.User{
		Email:     command.Email,
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Password:  command.Password,
		Provider:  string(provider.LocalProvider),
	}
	userService := services.NewUserService(api.store)
	if err := userService.CreateUser(user); err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			return response.Error(http.StatusConflict, "", err)
		} else {
			return response.Error(http.StatusInternalServerError, "", err)
		}
	}
	return response.JSON(http.StatusOK, map[string]string{
		"message": "USER_SUCCESSFULLY_SIGNED_UP",
	})
}

func (api *API) Reset(r *http.Request) response.Response {
	return response.JSON(http.StatusOK, map[string]string{
		"message": "PASSWORD_RESET_SUCCESSFULLY",
	})
}
