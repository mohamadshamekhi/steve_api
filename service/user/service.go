package userservice

import (
	"context"
	"github.com/abolfazlnorzad/steve/entity"
)

type Config struct {
	Prefix string `koanf:"prefix"`
}
type Repository interface {
	GetUserByPhoneNumber(ctx context.Context, phone string) (entity.User, error)
}
