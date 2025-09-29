package aclhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	permissionservice "github.com/abolfazlnorzad/steve/service/acl/permission"
	roleservice "github.com/abolfazlnorzad/steve/service/acl/role"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
)

type Handler struct {
	permissionSvc permissionservice.Service
	roleSvc       roleservice.Service
	translator    *i18n.Translator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	permissionSvc permissionservice.Service,
	roleSvc roleservice.Service,
	translator *i18n.Translator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,

) Handler {
	return Handler{
		permissionSvc: permissionSvc,
		roleSvc:       roleSvc,
		translator:    translator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	lang := header.Get("Accept-Language")
	if lang == "en" {
		return "en"
	}
	return "fa"
}
