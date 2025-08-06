package attributedefinitionhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	attributeservice "github.com/abolfazlnorzad/steve/service/attribute"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	attributedefinitionvalidator "github.com/abolfazlnorzad/steve/validator/attributedefinition"
)

type Handler struct {
	authSvc    authservice.Service
	authCfg    authservice.Config
	aclSvc     aclservice.Service
	translator *i18n.Translator
	vld        attributedefinitionvalidator.Validator
	svc        attributeservice.AttributeDefinition
}

func New(authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
	translator *i18n.Translator,
	vld attributedefinitionvalidator.Validator,
	svc attributeservice.AttributeDefinition) Handler {
	return Handler{
		authSvc:    authSvc,
		authCfg:    authCfg,
		aclSvc:     aclSvc,
		translator: translator,
		vld:        vld,
		svc:        svc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
