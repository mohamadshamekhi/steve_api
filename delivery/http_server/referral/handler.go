package referralhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
	referralvalidator "github.com/abolfazlnorzad/steve/validator/referral"
)

type Handler struct {
	svc        referralservice.StorefrontService
	validator  referralvalidator.Validator
	translator *i18n.Translator
	authSvc    authservice.Service
	authCfg    authservice.Config
}

func New(
	svc referralservice.StorefrontService,
	validator referralvalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
) Handler {
	return Handler{
		svc:        svc,
		validator:  validator,
		translator: translator,
		authSvc:    authSvc,
		authCfg:    authCfg,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
