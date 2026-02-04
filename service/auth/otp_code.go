package authservice

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

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
		otpCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		mErr := s.melliPayamakProvider.SendOTP(otpCtx, phoneNumber, code)
		if mErr != nil {
			dErr := s.repo.DeleteUserActiveCode(ctx, key)
			if dErr != nil {
				logger.L().Error("failed to delete user active code from db. key: ", key, "err: ", dErr)
			}
			logger.L().Error("failed to send otp: ", phoneNumber, "err: ", mErr)
			return richerror.New(op).WithWarpError(mErr)
		}
	}

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
