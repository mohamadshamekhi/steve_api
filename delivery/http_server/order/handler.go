package orderhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	ordervalidator "github.com/abolfazlnorzad/steve/validator/order"
)

type Handler struct {
	storefrontSvc orderservice.StoreFrontService
	vld           ordervalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	storefrontSvc orderservice.StoreFrontService,
	vld ordervalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service) Handler {
	return Handler{
		storefrontSvc: storefrontSvc,
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
