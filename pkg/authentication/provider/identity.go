package provider

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/square/go-jose/v3/jwt"
)

type AuthenticationProvider string

const (
	LDAPProvider  AuthenticationProvider = "ldap"
	LocalProvider AuthenticationProvider = "local"
)

type Identity struct {
	ID        *uuid.UUID
	Email     string
	FirstName string
	LastName  string
	Provider  AuthenticationProvider
	Username  string
}

func (identity *Identity) GetFullName() string {
	return fmt.Sprintf("%s %s", identity.FirstName, identity.LastName)
}

func (identity *Identity) ToClaims(attributes token.Attributes) token.Claims {
	return token.Claims{
		Claims: &jwt.Claims{
			Expiry:  jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			Subject: identity.GetFullName(),
		},
		User: token.User{
			ID:         identity.ID,
			Attributes: attributes,
			Audience:   "",
			Email:      identity.Email,
			FirstName:  identity.FirstName,
			LastName:   identity.LastName,
		},
	}
}
