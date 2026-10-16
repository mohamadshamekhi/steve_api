package referralhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
	referralvalidator "github.com/abolfazlnorzad/steve/validator/referral"
)

type Handler struct {
	svc           referralservice.StorefrontService
	backofficeSvc referralservice.BackofficeService
	validator     referralvalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	svc referralservice.StorefrontService,
	backofficeSvc referralservice.BackofficeService,
	validator referralvalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		svc:           svc,
		backofficeSvc: backofficeSvc,
		validator:     validator,
		translator:    translator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
