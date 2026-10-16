package orderdeliveryhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	orderdeliveryservice "github.com/abolfazlnorzad/steve/service/orderdelivery"
	orderdeliveryvalidator "github.com/abolfazlnorzad/steve/validator/orderdelivery"
)

type Handler struct {
	svc        orderdeliveryservice.Service
	vld        orderdeliveryvalidator.Validator
	translator *i18n.Translator
	authSvc    authservice.Service
	authCfg    authservice.Config
	aclSvc     aclservice.Service
}

func New(
	svc orderdeliveryservice.Service,
	vld orderdeliveryvalidator.Validator,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service) Handler {
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
