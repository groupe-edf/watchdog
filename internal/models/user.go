package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("INVALID_USERNAME_OR_PASSWORD")
	ErrUserNotFound              = errors.New("USER_NOT_FOUND")
	ErrUserAlreadyExists         = errors.New("USER_ALREADY_EXISTS")
)

// User represents a user
// swagger:model
type User struct {
	// The user's id
	ID             *uuid.UUID `json:"id"`
	ChangePassword bool       `json:"-"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	// swagger:strfmt email
	Email              string     `json:"email,omitempty"`
	FailedAttempts     int        `json:"-"`
	FirstName          string     `json:"first_name,omitempty"`
	LastLogin          *time.Time `json:"last_login,omitempty"`
	LastName           string     `json:"last_name,omitempty"`
	Locked             bool       `json:"locked"`
	LockedAt           *time.Time `json:"locked_at,omitempty"`
	MustChangePassword bool       `json:"must_change_password"`
	State              bool       `json:"state"`
	Password           string     `json:"password"`
	PasswordResetToken string     `json:"-"`
	Provider           string     `json:"provider"`
	UpdatedBy          *uuid.UUID `json:"-"`
	UpdatedAt          *time.Time `json:"-"`
	// The user's username
	Username string `json:"username,omitempty"`
}

// SetLastLogin set time to last login
func (user *User) SetLastLogin() {
	timestamp := time.Now()
	user.LastLogin = &timestamp
}

func (user *User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
