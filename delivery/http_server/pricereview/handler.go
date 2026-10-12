package pricereviewhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	pricereviewservice "github.com/abolfazlnorzad/steve/service/pricereview"
)

type Handler struct {
	backofficeSvc pricereviewservice.BackofficeService
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	backofficeSvc pricereviewservice.BackofficeService,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
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
