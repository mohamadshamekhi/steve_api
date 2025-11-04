package contactushandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	contactusservice "github.com/abolfazlnorzad/steve/service/contactus"
	contactusvalidator "github.com/abolfazlnorzad/steve/validator/contactus"
)

type Handler struct {
	backofficeSvc contactusservice.BackofficeService
	storefrontSvc contactusservice.StorefrontService
	translator    *i18n.Translator
	validator     contactusvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}

func New(
	storefrontSvc contactusservice.StorefrontService,
	translator *i18n.Translator,
	validator contactusvalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
	backofficeSvc contactusservice.BackofficeService,
) Handler {
	return Handler{
		translator:    translator,
		validator:     validator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
		storefrontSvc: storefrontSvc,
		backofficeSvc: backofficeSvc,
	}
}
