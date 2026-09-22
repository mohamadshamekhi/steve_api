package dashboardhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	dashboardservice "github.com/abolfazlnorzad/steve/service/dashboard"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc        dashboardservice.Service
	translator *i18n.Translator
	authSvc    authservice.Service
	authCfg    authservice.Config
}

func New(svc dashboardservice.Service, translator *i18n.Translator, authSvc authservice.Service, authCfg authservice.Config) Handler {
	return Handler{svc: svc, translator: translator, authSvc: authSvc, authCfg: authCfg}
}

func (h Handler) success(c echo.Context, data any) error {
	lang := "fa"
	if c.Request().Header.Get("Accept-Language") == "en" {
		lang = "en"
	}
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, data))
}
