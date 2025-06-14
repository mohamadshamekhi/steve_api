package backofficeauthhandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
	"net/http"
)

type Handler struct {
	svc        backofficeservice.Service
	vld        backofficevalidatorauth.Validator
	translator *i18n.Translator
}

func New(svc backofficeservice.Service, vld backofficevalidatorauth.Validator, translator *i18n.Translator) Handler {
	return Handler{
		svc:        svc,
		vld:        vld,
		translator: translator,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa" // زبان پیش‌فرض
}
