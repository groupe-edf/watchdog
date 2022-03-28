package middleware

import (
	"errors"
	"net/http"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/services"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

var ErrNotAuthorized error = errors.New("NOT_AUTHORIZED")

// Authenticator authentication middleware
type Authenticator struct {
	Logger                logging.Interface
	RequireAuthentication bool
	TokenService          token.TokenService
}

// Wrap implements Middleware
func (middleware *Authenticator) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleware.RequireAuthentication {
			user, err := middleware.checkAuthentication(r)
			if err != nil {
				data := response.Error(http.StatusUnauthorized, "NOT_AUTHORIZED", err)
				w.WriteHeader(data.Status())
				w.Write(data.Body())
				return
			}
			// Populate user info to request context
			r = r.Clone(token.SetUser(r.Context(), user))
		}
		// Update the current request with the new context information.
		next.ServeHTTP(w, r)
	})
}

func (middleware *Authenticator) checkAuthentication(r *http.Request) (token.User, error) {
	// APIKey authentication check
	tokenString := r.Header.Get(models.AccessTokenHeader)
	if tokenString != "" {
		di := container.GetContainer()
		store := di.Get(store.ServiceName).(models.Store)
		accessToken, err := store.FindAccessToken(tokenString)
		if err != nil {
			return token.User{}, err
		}
		if accessToken.IsExpired() || accessToken.Revoked {
			return token.User{}, models.ErrInvalidAccessToken
		}
		return token.User{}, nil
	}
	// Basic authentication check
	email, password, ok := r.BasicAuth()
	if ok {
		identity, err := services.Authenticate(email, password)
		if err != nil {
			return token.User{}, err
		}
		return identity.ToClaims(token.Attributes{}).User, nil
	}
	// JWT authentication check
	userToken, err := middleware.TokenService.Get(r)
	if err != nil {
		return token.User{}, err
	}
	if userToken != "" {
		claims, err := middleware.TokenService.Parse(userToken)
		if err != nil {
			return token.User{}, err
		}
		return claims.User, nil
	}
	return token.User{}, ErrNotAuthorized
}

// NewAuthenticator return authentication middleware
func NewAuthenticator(options token.JWTOptions) *Authenticator {
	logger := container.GetContainer().Get(logging.ServiceName).(logging.Interface)
	tokenService := token.NewJWTService(options)
	return &Authenticator{
		Logger:       logger,
		TokenService: tokenService,
	}
}
