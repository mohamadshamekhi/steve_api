package commenthandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	commentservice "github.com/abolfazlnorzad/steve/service/comment"
	commentvalidator "github.com/abolfazlnorzad/steve/validator/comment"
	"net/http"
)

type Handler struct {
	svc        commentservice.Service
	translator *i18n.Translator
	validator  commentvalidator.Validator
	authSvc    authservice.Service
	authCfg    authservice.Config
	aclSvc     aclservice.Service
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}

func New(svc commentservice.Service,
	translator *i18n.Translator,
	validator commentvalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service) Handler {
	return Handler{
		svc:        svc,
		translator: translator,
		validator:  validator,
		authSvc:    authSvc,
		authCfg:    authCfg,
		aclSvc:     aclSvc,
	}
}
