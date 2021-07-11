package token

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrParsingUser error when parsing user from context
	ErrParsingUser = errors.New("User can't be parsed")
)

// ContextKey is used for context.Context value. The value requires a key that is not primitive type.
type ContextKey string

const (
	// UserContextKey is the ContextKey for User
	UserContextKey ContextKey = "user"
)

// User is the basic part of oauth data provided by service
type User struct {
	ID         *uuid.UUID `json:"id"`
	Attributes Attributes `json:"attributes,omitempty"`
	Audience   string     `json:"audience,omitempty"`
	Email      string     `json:"email,omitempty"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
}

// Attributes user custom attributes
type Attributes struct {
	items map[string]interface{}
}

// AddAttribute add attribute
func (attributes *Attributes) AddAttribute(key string, value interface{}) {
	attributes.items[key] = value
}

// BoolAttribute get boolean attribute value
func (attributes *Attributes) BoolAttribute(key string) bool {
	result, ok := attributes.items[key].(bool)
	if !ok {
		return false
	}
	return result
}

// GetAttribute get attribute value
func (attributes *Attributes) GetAttribute(key string) interface{} {
	if attributes.items[key] != nil {
		return attributes.items[key]
	}
	return nil
}

// StringAttribute get string attribute value
func (attributes *Attributes) StringAttribute(key string) string {
	result, ok := attributes.GetAttribute(key).(string)
	if !ok {
		return ""
	}
	return result
}

// GetUser returns user info from context
func GetUser(ctx context.Context) (user User, err error) {
	if ctx == nil {
		return user, errors.New("no user data found")
	}
	rawUser := ctx.Value(UserContextKey)
	if user, ok := rawUser.(User); ok {
		return user, nil
	}
	return user, ErrParsingUser
}

// SetUser sets user into request context
func SetUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}
