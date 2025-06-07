package backofficeservice

import (
	"context"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
)

type Config struct {
	Prefix string `koanf:"prefix"`
}
type Repository interface {
	GetBackOfficeUserIDByPhoneNumber(ctx context.Context, phoneNumber string) (string, error)
}
type Service struct {
	config  Config
	authSvc authservice.Service
	repo    Repository
}

func New(cfg Config, authSvc authservice.Service, repo Repository) Service {
	return Service{
		config:  cfg,
		authSvc: authSvc,
		repo:    repo,
	}
}
