package models

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/google/uuid"
)

const AccessTokenHeader = "X-API-Key"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
)

// AccessToken represents a user access token
// swagger:model
type AccessToken struct {
	ID         int64      `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Name       string     `json:"name"`
	Revoked    bool       `json:"revoked"`
	Scopes     []string   `json:"scopes"`
	Token      string     `json:"token"`
	UserID     *uuid.UUID `json:"user_id"`
}

func (token *AccessToken) IsExpired() bool {
	today := time.Now()
	if token.ExpiresAt != nil && today.Before(*token.ExpiresAt) {
		return true
	}
	return false
}

func (token *AccessToken) Revoke() {
	token.Revoked = true
}

func NewAccessToken(name string, userID *uuid.UUID) *AccessToken {
	createdAt := time.Now()
	return &AccessToken{
		CreatedAt: createdAt,
		Name:      name,
		Token:     GenerateToken(32),
		UserID:    userID,
	}
}

func GenerateToken(length int) string {
	const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for index, b := range bytes {
		bytes[index] = characters[b%byte(len(characters))]
	}
	return string(bytes)
}
