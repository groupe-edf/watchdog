package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/authentication/encoder"
)

type UserService struct {
	store models.Store
}

func (service *UserService) CreateUser(user *models.User) error {
	existingUser, err := service.store.FindUserByEmail(user.Email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return err
	}
	if existingUser != nil {
		return models.ErrUserAlreadyExists
	}
	encoder := encoder.NewBCryptEncoder()
	hashedPassword, err := encoder.HashPassword(user.Password)
	if err != nil {
		return err
	}
	id := uuid.New()
	user.ID = &id
	user.CreatedAt = time.Now()
	user.Password = string(hashedPassword)
	_, err = service.store.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(store models.Store) *UserService {
	return &UserService{
		store: store,
	}
}
