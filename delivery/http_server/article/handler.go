package articlehandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	articleservice "github.com/abolfazlnorzad/steve/service/article"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	articlevalidator "github.com/abolfazlnorzad/steve/validator/article"
	"net/http"
)

type Handler struct {
	translator           *i18n.Translator
	authSvc              authservice.Service
	authCfg              authservice.Config
	aclSvc               aclservice.Service
	backofficeArticleSvc articleservice.BackofficeService
	vld                  articlevalidator.Validator
}

func New(
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
	backofficeArticleSvc articleservice.BackofficeService,
	vld articlevalidator.Validator,
) Handler {
	return Handler{
		translator:           translator,
		authSvc:              authSvc,
		authCfg:              authCfg,
		aclSvc:               aclSvc,
		backofficeArticleSvc: backofficeArticleSvc,
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
