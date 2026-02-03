package authservice

import (
	"context"
	"time"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessSubject         string        `koanf:"access_subject"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	ExpirationTimeOptCode time.Duration `koanf:"expiration_time_opt_code"`
}

type Repository interface {
	CheckUserHasActiveCode(ctx context.Context, key string) (bool, error)
	GetUserActiveCode(ctx context.Context, key string) (string, error)
	StoreNewActiveCodeForUser(ctx context.Context, key string, code string, expireTime time.Duration) error
	IsInBlocklist(ctx context.Context, kind string, JWTID string) (bool, error)
	DeleteUserActiveCode(ctx context.Context, key string) error
}

type MelliPayamakProvider interface {
	SendOTP(ctx context.Context, to string, code string) error
}
type Service struct {
	cfg                  Config
	repo                 Repository
	melliPayamakProvider MelliPayamakProvider
}

func New(cfg Config, repo Repository, melliPayamakProvider MelliPayamakProvider) Service {
	return Service{cfg: cfg, repo: repo, melliPayamakProvider: melliPayamakProvider}
}

func (s Service) CheckTokenIsInBlockList(ctx context.Context, kind string, id string) (bool, error) {
	return s.repo.IsInBlocklist(ctx, kind, id)
}
