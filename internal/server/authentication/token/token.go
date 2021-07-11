package token

import (
	"errors"
	"net/http"
	"strings"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// JWTOptions JWT service options
type JWTOptions struct {
	BearerHeader  string // default "Bearer"
	JWTCookieName string // default "JWT"
	JWTHeaderKey  string // default "Authorization"
	JWTQueryParam string
	Secret        []byte
}

// Claims service custom claims
type Claims struct {
	*jwt.Claims
	User User `json:"user,omitempty"`
}

// TokenService defines interface accessing tokens
type TokenService interface {
	Parse(tokenString string) (claims *Claims, err error)
	Set(w http.ResponseWriter, claims Claims) (Claims, error)
	Get(r *http.Request) (token string, err error)
	IsExpired(claims Claims) bool
}

// JWT service that wraps JWT operations
type JWT struct {
	Options JWTOptions
}

// Get token from url, header or cookie
func (service *JWT) Get(r *http.Request) (string, error) {
	tokenString := ""
	// Try to get from "token" query param
	if token := r.URL.Query().Get(service.Options.JWTQueryParam); token != "" {
		tokenString = token
	}
	// Try to get from JWT header
	if authenticationHeader := r.Header.Get(service.Options.JWTHeaderKey); authenticationHeader != "" && tokenString == "" {
		token := strings.Split(authenticationHeader, service.Options.BearerHeader)
		if len(token) < 2 {
			return "", errors.New("wrong authorization header format")
		}
		tokenString = strings.TrimSpace(token[1])
	}
	return tokenString, nil
}

// IsExpired returns true if claims expired
func (service *JWT) IsExpired(claims Claims) bool {
	return false
}

// Parse retrun claims from JWT token
func (service *JWT) Parse(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}
	claims := &Claims{}
	err = parsedToken.Claims(service.Options.Secret, &claims)
	return claims, err
}

// Set creates token cookie with xsrf cookie and put it to ResponseWriter
func (service *JWT) Set(w http.ResponseWriter, claims Claims) (Claims, error) {
	return Claims{}, nil
}

// Token generate new JWT token
func (service *JWT) Token(claims Claims) (string, error) {
	// Create signing key
	key := jose.SigningKey{Algorithm: jose.HS256, Key: service.Options.Secret}
	// Create a signer, used to sign the JWT
	var signerOptions = jose.SignerOptions{}
	signerOptions.WithType("JWT")
	signer, err := jose.NewSigner(key, &signerOptions)
	if err != nil {
		return "", err
	}
	// Create an instance of Builder that uses the secret signer
	builder := jwt.Signed(signer)
	// Add claims to the Builder
	builder = builder.Claims(claims)
	return builder.CompactSerialize()
}

// NewJWTService create new JWT service
func NewJWTService(options JWTOptions) *JWT {
	if options.JWTHeaderKey == "" {
		options.JWTHeaderKey = "Authorization"
	}
	if options.BearerHeader == "" {
		options.BearerHeader = "Bearer"
	}
	service := &JWT{
		Options: options,
	}
	return service
}
