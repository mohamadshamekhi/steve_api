package sellhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	sellservice "github.com/abolfazlnorzad/steve/service/sell"
	sellvalidator "github.com/abolfazlnorzad/steve/validator/sell"
)

type Handler struct {
	backofficeSvc sellservice.BackofficeService
	storefrontSvc sellservice.StorefrontService
	translator    *i18n.Translator
	validator     sellvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(backofficeSvc sellservice.BackofficeService,
	storefrontSvc sellservice.StorefrontService,
	translator *i18n.Translator,
	validator sellvalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		storefrontSvc: storefrontSvc,
		translator:    translator,
		validator:     validator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
