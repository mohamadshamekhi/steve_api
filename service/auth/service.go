package authservice

import (
	"context"
	"fmt"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	"github.com/abolfazlnorzad/steve/contract/goproto/auth"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"math/rand/v2"
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

func (s Service) GenerateOTPCode(ctx context.Context, phoneNumber string, prefix string) error {
	const op = "authservice.GenerateOTPCode"
	key := fmt.Sprintf("%s:%s", prefix, phoneNumber)
	resp, cErr := s.repo.CheckUserHasActiveCode(ctx, key)
	if cErr != nil {
		return richerror.New(op).WithWarpError(cErr)
	}
	if !resp {
		code := fmt.Sprintf("%d", 10000+rand.IntN(99999))
		sErr := s.repo.StoreNewActiveCodeForUser(ctx, key, code, s.cfg.ExpirationTimeOptCode)
		if sErr != nil {
			return richerror.New(op).WithWarpError(sErr)
		}
		event := &auth.SendOTPCode{Code: code, PhoneNumber: phoneNumber}
		payload, mErr := proto.Marshal(event)
		if mErr != nil {
			return richerror.New(op).WithWarpError(mErr)
		}
		fmt.Println("ppppppp  :", payload)
		s.nats.NatsJS().AddStream(&nats.StreamConfig{
			Name:     "notification",
			MaxAge:   1 * 24 * time.Hour, // حداکثر عمر پیام (مثلاً 1 روز)
			Subjects: []string{"notification.*"},
		})
		// publish code in nats (phone_number , code)
		ack, pErr := s.nats.NatsJS().Publish("notification.send_otp", payload)
		if pErr != nil {
			fmt.Println("eeeeeeeeeeeee", pErr)
			return richerror.New(op).WithWarpError(pErr)
		}
		logger.L().Info("message published successfully, ack.Stream , ack.Sequence", ack.Stream, ack.Sequence)
	}

	return nil
}
