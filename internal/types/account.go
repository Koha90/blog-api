package types

import (
	"math/rand"
	"time"
)

type Account struct {
	ID                int
	FirstName         string
	LastName          string
	Number            int64
	EncryptedPassword string
	Balance           float64
	CreatedAt         time.Time
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: password,
		Number:            int64(rand.Intn(100000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
