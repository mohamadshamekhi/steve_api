package backofficeauthhandler

import (
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
)

type Handler struct {
	svc backofficeservice.Service
	vld backofficevalidatorauth.Validator
}

func New(svc backofficeservice.Service, vld backofficevalidatorauth.Validator) Handler {
	return Handler{
		svc: svc,
		vld: vld,
	}
}
