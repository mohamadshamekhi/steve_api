package backofficeservice

import (
	"context"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) Login(ctx context.Context, req backofficeauthparam.SendLoginCodeRequest) (backofficeauthparam.SendLoginCodeResponse, error) {
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
