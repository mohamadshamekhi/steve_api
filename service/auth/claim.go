package authservice

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	EntityID      string   `json:"entity_id"`
	Kind          string   `json:"kind"`
	RoleIDs       []string `json:"role_ids"`
	PermissionIDs []string `json:"permission_ids"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
