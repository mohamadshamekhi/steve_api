package supplierhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/pkg/i18n"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	supplierservice "github.com/abolfazlnorzad/steve/service/supplier"
	suppliervalidator "github.com/abolfazlnorzad/steve/validator/supplier"
	"github.com/labstack/echo/v4"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
)

type Handler struct {
	backofficeSvc supplierservice.BackofficeService
	translator    *i18n.Translator
	validator     suppliervalidator.Validator
	authSvc       authservice.Service
	authCfg       authservice.Config
	aclSvc        aclservice.Service
}

func New(
	backofficeSvc supplierservice.BackofficeService,
	translator *i18n.Translator,
	validator suppliervalidator.Validator,
	authSvc authservice.Service,
	authCfg authservice.Config,
	aclSvc aclservice.Service,
) Handler {
	return Handler{
		backofficeSvc: backofficeSvc,
		translator:    translator,
		validator:     validator,
		authSvc:       authSvc,
		authCfg:       authCfg,
		aclSvc:        aclSvc,
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}

func (h Handler) success(c echo.Context, data any) error {
	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, data))
}
