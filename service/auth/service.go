package authservice

import (
	"context"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
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
}

type Service struct {
	cfg  Config
	repo Repository
	nats *natsadapter.Adapter
}

func New(cfg Config, repo Repository, nats *natsadapter.Adapter) Service {
	return Service{cfg: cfg, repo: repo, nats: nats}
}
