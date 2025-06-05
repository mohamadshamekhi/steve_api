package backofficeservice

import (
	"context"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) SendOTPCode(ctx context.Context, req backofficeauthparam.SendLoginCodeRequest) (backofficeauthparam.SendLoginCodeResponse, error) {
	const op = "backofficeservice.Login"

	err := s.authSvc.GenerateOTPCode(ctx, req.PhoneNumber, s.config.Prefix)
	if err != nil {
		return backofficeauthparam.SendLoginCodeResponse{}, richerror.New(op)
	}

	// return resp.
	return backofficeauthparam.SendLoginCodeResponse{
		Message: "sms is send",
	}, nil
}

func (s Service) VerifyOTPCode(ctx context.Context, req backofficeauthparam.VerifyLoginCodeRequest) (backofficeauthparam.VerifyLoginCodeResponse, error) {
	const op = "backofficeservice.VerifyOTPCode"

	resp, err := s.authSvc.VerifyOTPCode(ctx, req.PhoneNumber, req.Code, s.config.Prefix)
	if err != nil {
		return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op)
	}
	if resp {
		token, gErr := s.authSvc.GenerateToken(req.PhoneNumber)
		return backofficeauthparam.VerifyLoginCodeResponse{
			Token:   "",
			Message: "you are logged in!",
		}, nil
	}
	return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).
		WithMessage(errmsg.ErrorOTPCodeIsNotCorrect)
}
