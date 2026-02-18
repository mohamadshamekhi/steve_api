package searchhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	searchservice "github.com/abolfazlnorzad/steve/service/search"
)

type Handler struct {
	translator          *i18n.Translator
	storefrontSearchSvc searchservice.StorefrontService
}

func New(
	translator *i18n.Translator,
	storefrontSearchSvc searchservice.StorefrontService,
) Handler {
	return Handler{
		translator:          translator,
		storefrontSearchSvc: storefrontSearchSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
