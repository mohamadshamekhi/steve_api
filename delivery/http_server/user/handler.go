package userhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	uservalidator "github.com/abolfazlnorzad/steve/validator/user"
)

type Handler struct {
	storeFrontSvc userservice.StoreFrontService
	backofficeSvc userservice.BackofficeService
	vld           uservalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	acl           aclservice.Service
}

func New(backofficeSvc userservice.BackofficeService, storeFrontSvc userservice.StoreFrontService, vld uservalidator.Validator, translator *i18n.Translator,
	authSvc authservice.Service, authCfg authservice.Config, acl aclservice.Service) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		storeFrontSvc: storeFrontSvc,
		vld:           vld,
		translator:    translator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		acl:           acl,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
