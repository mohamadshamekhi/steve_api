package brandhandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
	"net/http"
)

type Handler struct {
	backofficeSvc brandservice.BackofficeService
	translator    *i18n.Translator
	validator     brandvalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
}

func New(backofficeSvc brandservice.BackofficeService,
	translator *i18n.Translator,
	validator brandvalidator.Validator, authSvc authservice.Service, authCfg authservice.Config) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		translator:    translator,
		validator:     validator,
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
