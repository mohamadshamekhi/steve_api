package backofficehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	backofficevalidator "github.com/abolfazlnorzad/steve/validator/backoffice"
)

type Handler struct {
	svc        backofficeservice.Service
	vld        backofficevalidator.Validator
	translator *i18n.Translator
	authSvc    authservice.Service
	authCfg    authservice.Config
	aclSvc     aclservice.Service
}

func New(svc backofficeservice.Service, vld backofficevalidator.Validator, translator *i18n.Translator,
	authSvc authservice.Service, authCfg authservice.Config, aclSvc aclservice.Service) Handler {
	return Handler{
		svc:        svc,
		vld:        vld,
		translator: translator,
		authSvc:    authSvc,
		authCfg:    authCfg,
		aclSvc:     aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
