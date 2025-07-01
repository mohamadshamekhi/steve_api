package mediahandler

import (
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
	"net/http"
)

type Handler struct {
	vld        mediavalidator.Validator
	svc        mediaservice.Service
	translator *i18n.Translator
	authSvc    authservice.Service
	authCfg    authservice.Config
}

func New(vld mediavalidator.Validator,
	svc mediaservice.Service,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config) Handler {
	return Handler{
		vld:        vld,
		svc:        svc,
		translator: translator,
		authSvc:    authSvc,
		authCfg:    authCfg,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
