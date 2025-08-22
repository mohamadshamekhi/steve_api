package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	cartservice "github.com/abolfazlnorzad/steve/service/cart"
	cartvalidator "github.com/abolfazlnorzad/steve/validator/cart"
)

type Handler struct {
	storefrontSvc cartservice.StoreFrontService
	backofficeSvc cartservice.BackofficeService
	vld           cartvalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(storefrontSvc cartservice.StoreFrontService,
	backofficeSvc cartservice.BackofficeService,
	vld cartvalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,

	aclSvc aclservice.Service) Handler {
	return Handler{
		storefrontSvc: storefrontSvc,
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
