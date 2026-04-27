package authservice

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
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
		code, gErr := s.generate5DigitOTP()
		if gErr != nil {
			return richerror.New(op).WithWarpError(gErr)
		}
		sErr := s.repo.StoreNewActiveCodeForUser(ctx, key, code, s.cfg.ExpirationTimeOptCode)
		if sErr != nil {
			return richerror.New(op).WithWarpError(sErr)
		}
		otpCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ExpirationTimeOtpContext)
		defer cancel()
		mErr := s.melliPayamakProvider.SendPatternSMS(otpCtx, phoneNumber, []string{code}, mellipayamak.PatternOTPKey)
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
func (s Service) generate5DigitOTP() (string, error) {
	maxNumber := big.NewInt(90000)

	n, err := rand.Int(rand.Reader, maxNumber)
	if err != nil {
		return "", err
	}

	secureCode := n.Int64() + 10000

	return fmt.Sprintf("%d", secureCode), nil
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
