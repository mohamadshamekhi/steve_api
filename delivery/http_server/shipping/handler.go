package shippinghandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	settingservice "github.com/abolfazlnorzad/steve/service/setting"
	settingvalidator "github.com/abolfazlnorzad/steve/validator/setting"
)

type Handler struct {
	backofficeSvc settingservice.BackofficeShippingService
	vld           settingvalidator.ShippingValidator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	backofficeSvc settingservice.BackofficeShippingService,
	vld settingvalidator.ShippingValidator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		vld:           vld,
		translator:    translator,
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
