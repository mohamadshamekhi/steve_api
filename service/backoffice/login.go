package backofficeservice

import (
	"fmt"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) Login(req backofficeauthparam.SendLoginCodeRequest) (backofficeauthparam.SendLoginCodeResponse, error) {
	const op = "backofficeservice.Login"

	code, err := s.authSvc.GenerateOTPCode(req.PhoneNumber, s.config.Prefix)
	if err != nil {
		return backofficeauthparam.SendLoginCodeResponse{}, richerror.New(op)
	}
	// publish message into nats phone_number and code with proto_buf serialization
	fmt.Println("code : ", code)
	// publish code to nats to consume this code with notification service

	// return resp.
	return backofficeauthparam.SendLoginCodeResponse{
		Message: "sms is send",
	}, nil
}
