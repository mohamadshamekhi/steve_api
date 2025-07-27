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
	GetUserByID(ctx context.Context, phone string) (entity.User, error)
	AddNewUser(ctx context.Context, u entity.User) error
	GetOrCreateUser(ctx context.Context, phonenumber string) (entity.User, error)
	UpdateUser(ctx context.Context, u entity.User) error
}
