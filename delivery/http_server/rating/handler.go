package ratinghandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	ratingservice "github.com/abolfazlnorzad/steve/service/rating"
	ratingvalidator "github.com/abolfazlnorzad/steve/validator/rating"
)

type Handler struct {
	service    ratingservice.Service
	translator *i18n.Translator
	validator  ratingvalidator.Validator
	authSvc    authservice.Service
	authCfg    authservice.Config
	aclSvc     aclservice.Service
}

func New(
	service ratingservice.Service,
	translator *i18n.Translator,
	validator ratingvalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service) Handler {
	return Handler{
		service:    service,
		translator: translator,
		validator:  validator,
		authSvc:    authSvc,
		authCfg:    authCfg,
		aclSvc:     aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
