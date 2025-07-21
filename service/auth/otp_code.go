package authservice

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/contract/goproto/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"time"
)

func (s Service) GenerateOTPCode(ctx context.Context, phoneNumber string, prefix string) error {
	const op = "authservice.GenerateOTPCode"
	key := fmt.Sprintf("%s:%s", prefix, phoneNumber)
	resp, cErr := s.repo.CheckUserHasActiveCode(ctx, key)
	if cErr != nil {
		return richerror.New(op).WithWarpError(cErr)
	}
	if !resp {
		// todo - uncomment !!! important !!!
		//code := fmt.Sprintf("%d", 10000+rand.IntN(99999))
		code := "77777"
		sErr := s.repo.StoreNewActiveCodeForUser(ctx, key, code, s.cfg.ExpirationTimeOptCode)
		if sErr != nil {
			return richerror.New(op).WithWarpError(sErr)
		}
		return nil
		//return s.publishMessageToNats(phoneNumber, code)
	}

	return nil
}

func (s Service) publishMessageToNats(phoneNumber string, code string) error {
	const op = "authservice.GenerateOTPCode"

	event := &auth.SendOTPCode{Code: code, PhoneNumber: phoneNumber}
	payload, mErr := proto.Marshal(event)
	if mErr != nil {
		return richerror.New(op).WithWarpError(mErr).WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	// todo - remove this code and transfer to package that init all stream or use Terraform or k8s opt
	_, aErr := s.nats.NatsJS().AddStream(&nats.StreamConfig{
		Name:     "notification",
		MaxAge:   1 * 24 * time.Hour,
		Subjects: []string{"notification.*"},
	})
	if aErr != nil {
		return richerror.New(op).WithWarpError(aErr).WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	ack, pErr := s.nats.NatsJS().Publish(constant.SendOTPSubject, payload)
	if pErr != nil {
		return richerror.New(op).WithWarpError(pErr).WithKind(richerror.KindUnexpected).WithMessageKey(errmsg.ErrorMsgSomethingWentWrong)
	}
	// todo - if necessary, we can store this published message into db.
	logger.L().Info("message published successfully, ack.Stream , ack.Sequence", ack.Stream, ack.Sequence)
	return nil
}

func (s Service) VerifyOTPCode(ctx context.Context, phoneNumber string, code string, prefix string) (bool, error) {
	const op = "authservice.VerifyOTPCode"
	key := fmt.Sprintf("%s:%s", prefix, phoneNumber)
	activeCode, gErr := s.repo.GetUserActiveCode(ctx, key)
	if gErr != nil {
		return false, richerror.New(op).WithMessageKey("validation.incorrect_opt_code").WithWarpError(gErr)
	}
	if activeCode != code {
		return false, nil
	}
	return true, nil
}
