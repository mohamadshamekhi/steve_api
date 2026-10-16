package partnerhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	partnervalidator "github.com/abolfazlnorzad/steve/validator/partner"
)

type Handler struct {
	storefrontSvc partnerservice.StorefrontService
	shopSvc       partnerservice.ShopService
	backofficeSvc partnerservice.BackofficeService
	validator     partnervalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	storefrontSvc partnerservice.StorefrontService,
	shopSvc partnerservice.ShopService,
	backofficeSvc partnerservice.BackofficeService,
	validator partnervalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		storefrontSvc: storefrontSvc,
		shopSvc:       shopSvc,
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
