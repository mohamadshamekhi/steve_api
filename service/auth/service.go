package authservice

import (
	"context"
	"fmt"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	"github.com/abolfazlnorzad/steve/contract/goproto/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
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
			return richerror.New(op).WithWarpError(mErr).WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrorMsgSomethingWentWrong)
		}
		// todo - remove this code and transfer to package that init all stream or use Terraform or k8s opt
		_, aErr := s.nats.NatsJS().AddStream(&nats.StreamConfig{
			Name:     "notification",
			MaxAge:   1 * 24 * time.Hour,
			Subjects: []string{"notification.*"},
		})
		if aErr != nil {
			return richerror.New(op).WithWarpError(aErr).WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrorMsgSomethingWentWrong)
		}
		ack, pErr := s.nats.NatsJS().Publish(constant.SendOTPSubject, payload)
		if pErr != nil {
			fmt.Println("pErr", pErr)
			return richerror.New(op).WithWarpError(pErr).WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrorMsgSomethingWentWrong)
		}
		logger.L().Info("message published successfully, ack.Stream , ack.Sequence", ack.Stream, ack.Sequence)
	}

	return nil
}
