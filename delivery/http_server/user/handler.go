package userhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	uservalidator "github.com/abolfazlnorzad/steve/validator/user"
)

type Handler struct {
	storeFrontSvc userservice.StoreFrontService
	vld           uservalidator.Validator
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
}

func New(storeFrontSvc userservice.StoreFrontService, vld uservalidator.Validator, translator *i18n.Translator,
	authSvc authservice.Service, authCfg authservice.Config) Handler {
	return Handler{
		storeFrontSvc: storeFrontSvc,
		vld:           vld,
		translator:    translator,
		authSvc:       authSvc,
		authCfg:       authCfg,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
