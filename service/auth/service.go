package authservice

import (
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
	CheckUserHasActiveCode(key string) (bool, error)
	GetUserActiveCode(key string) (string, error)
	StoreNewActiveCodeForUser(key string, code string) error
}

type Service struct {
	cfg  Config
	repo Repository
}

func New(cfg Config, repo Repository) Service {
	return Service{cfg: cfg, repo: repo}
}

func (s Service) GenerateOTPCode(phoneNumber string, prefix string) (string, error) {
	const op = "authservice.GenerateOTPCode"
	key := fmt.Sprintf("%s:%s", prefix, phoneNumber)
	resp, cErr := s.repo.CheckUserHasActiveCode(key)
	if cErr != nil {
		return "", richerror.New(op).WithWarpError(cErr)
	}
	var code string
	if !resp {
		randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
		code = fmt.Sprintf("%v", randomInt)
		sErr := s.repo.StoreNewActiveCodeForUser(key, code)
		if sErr != nil {
			return "", richerror.New(op).WithWarpError(sErr)
		}
	}
	code, gErr := s.repo.GetUserActiveCode(key)
	if gErr != nil {
		return "", richerror.New(op).WithWarpError(gErr)
	}

	return code, nil
}
