package producthandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	productvalidator "github.com/abolfazlnorzad/steve/validator/product"
	"net/http"
)

type Handler struct {
	translator           *i18n.Translator
	authSvc              authservice.Service
	authCfg              authservice.Config
	aclSvc               aclservice.Service
	backofficeProductSvc productservice.BackofficeService
	vld                  productvalidator.Validator
}

func NewHandler(
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
	backofficeProductSvc productservice.BackofficeService,
	vld productvalidator.Validator,
) Handler {
	return Handler{
		translator:           translator,
		authSvc:              authSvc,
		authCfg:              authCfg,
		aclSvc:               aclSvc,
		backofficeProductSvc: backofficeProductSvc,
		vld:                  vld,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
