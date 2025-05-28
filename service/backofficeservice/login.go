package backofficeservice

import "github.com/abolfazlnorzad/steve/param/backofficeparam"

func (s Service) Login(req backofficeparam.SendLoginCodeRequest) (backofficeparam.SendLoginCodeResponse, error) {
	// check phone number is exist

	// if exist -> connect to auth service and create code

	// (in auth service send code and user phone number into nats to send sms and return code to this service)

	// return resp.
}
