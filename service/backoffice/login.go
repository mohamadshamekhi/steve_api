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
		return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(err)
	}
	if resp {
		backofficeID, iErr := s.repo.GetBackOfficeUserIDByPhoneNumber(ctx, req.PhoneNumber)
		if iErr != nil {
			return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(iErr)
		}
		token, gErr := s.authSvc.GenerateAccessToken(backofficeID, s.config.Prefix)
		if gErr != nil {
			return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithWarpError(gErr).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}
		return backofficeauthparam.VerifyLoginCodeResponse{
			Token:   token,
			Message: "you are logged in!",
		}, nil
	}
	return backofficeauthparam.VerifyLoginCodeResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).
		WithMessage(errmsg.ErrorOTPCodeIsNotCorrect)
}
