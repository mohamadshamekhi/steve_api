package authservice

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	EntityID       string   `json:"entity_id"`
	Kind           string   `json:"kind"`
	RoleKeys       []string `json:"role_keys"`
	PermissionKeys []string `json:"permission_keys"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
