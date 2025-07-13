package categoryhandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	categoryservice "github.com/abolfazlnorzad/steve/service/category"
	categoryvalidator "github.com/abolfazlnorzad/steve/validator/category"
	"net/http"
)

type Handler struct {
	backofficeCategorySvc categoryservice.BackofficeService
	storefrontCategorySvc categoryservice.StorefrontService
	vld                   categoryvalidator.Validator
	translator            *i18n.Translator
	authSvc               authservice.Service
	authCfg               authservice.Config
	aclSvc                aclservice.Service
}

func New(
	backofficeCategorySvc categoryservice.BackofficeService,
	vld categoryvalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
	storefrontCategorySvc categoryservice.StorefrontService,
) Handler {
	return Handler{
		backofficeCategorySvc: backofficeCategorySvc,
		vld:                   vld,
		translator:            translator,
		authSvc:               authSvc,
		authCfg:               authCfg,
		aclSvc:                aclSvc,
		storefrontCategorySvc: storefrontCategorySvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
