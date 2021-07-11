package encoder

import "golang.org/x/crypto/bcrypt"

type Encoder interface {
	HashPassword(password string) ([]byte, error)
	IsPasswordValid(password string, plainPassword string) error
}

type BCryptEncoder struct{}

func (hasher *BCryptEncoder) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}

func (hasher *BCryptEncoder) IsPasswordValid(password string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(plainPassword))
}

func NewBCryptEncoder() Encoder {
	return &BCryptEncoder{}
}
