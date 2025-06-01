package authservice

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"math/rand"
	"time"
)

type Config struct {
	SignKey              string
	AccessSubject        string
	AccessExpirationTime time.Duration
}

type Repository interface {
	CheckUserHasActiveCode(ctx context.Context, key string) (bool, error)
	GetUserActiveCode(ctx context.Context, key string) (string, error)
	StoreNewActiveCodeForUser(ctx context.Context, key string, code string) error
}

type Service struct {
	cfg  Config
	repo Repository
}

func New(cfg Config, repo Repository) Service {
	return Service{cfg: cfg, repo: repo}
}

func (s Service) GenerateOTPCode(ctx context.Context, phoneNumber string, prefix string) (string, error) {
	const op = "authservice.GenerateOTPCode"
	key := fmt.Sprintf("%s:%s", prefix, phoneNumber)
	resp, cErr := s.repo.CheckUserHasActiveCode(ctx, key)
	if cErr != nil {
		return "", richerror.New(op).WithWarpError(cErr)
	}
	var code string
	if !resp {
		randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
		code = fmt.Sprintf("%v", randomInt)
		sErr := s.repo.StoreNewActiveCodeForUser(ctx, key, code)
		if sErr != nil {
			return "", richerror.New(op).WithWarpError(sErr)
		}
	}
	code, gErr := s.repo.GetUserActiveCode(ctx, key)
	if gErr != nil {
		return "", richerror.New(op).WithWarpError(gErr)
	}

	return code, nil
}
