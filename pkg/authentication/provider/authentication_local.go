package provider

import (
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/pkg/authentication/encoder"
	"github.com/groupe-edf/watchdog/pkg/container"
)

type Local struct {
}

func (provider *Local) Authenticate(email string, password string) (identity *Identity, err error) {
	store := container.Get(store.ServiceName).(models.Store)
	user, err := store.FindUserByEmail(email)
	if err != nil {
		return identity, err
	}
	hasher := encoder.NewBCryptEncoder()
	if err := hasher.IsPasswordValid(user.Password, password); err != nil {
		return identity, err
	}
	identity = &Identity{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Provider:  LocalProvider,
	}
	return identity, err
}

func NewLocalProvider() *Local {
	return &Local{}
}
