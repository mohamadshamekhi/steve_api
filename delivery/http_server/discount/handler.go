package discounthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	discountvalidator "github.com/abolfazlnorzad/steve/validator/discount"
)

type Handler struct {
	backofficeSvc discountservice.BackofficeService
	translator    *i18n.Translator
	validator     discountvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(backofficeSvc discountservice.BackofficeService,
	translator *i18n.Translator,
	validator discountvalidator.Validator, authSvc authservice.Service, authCfg authservice.Config, aclSvc aclservice.Service,
) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
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
