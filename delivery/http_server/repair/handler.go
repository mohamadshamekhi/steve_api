package repairhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	repairservice "github.com/abolfazlnorzad/steve/service/repair"
	repairvalidator "github.com/abolfazlnorzad/steve/validator/repair"
)

type Handler struct {
	backofficeSvc repairservice.BackofficeService
	storefrontSvc repairservice.StorefrontService
	translator    *i18n.Translator
	validator     repairvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(backofficeSvc repairservice.BackofficeService,
	storefrontSvc repairservice.StorefrontService,
	translator *i18n.Translator,
	validator repairvalidator.Validator,
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
