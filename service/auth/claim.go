package authservice

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	PhoneNumber string `json:"phone_number"`
	Kind        string `json:"kind"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
