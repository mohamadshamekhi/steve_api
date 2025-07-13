package brandhandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
	"net/http"
)

type Handler struct {
	backofficeSvc brandservice.BackofficeService
	storefrontSvc brandservice.StorefrontService
	translator    *i18n.Translator
	validator     brandvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(backofficeSvc brandservice.BackofficeService,
	translator *i18n.Translator,
	validator brandvalidator.Validator, authSvc authservice.Service, authCfg authservice.Config, aclSvc aclservice.Service,
	storefrontSvc brandservice.StorefrontService,

) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		translator:    translator,
		validator:     validator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
		storefrontSvc: storefrontSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
