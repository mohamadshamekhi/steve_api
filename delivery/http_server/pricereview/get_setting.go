package pricereviewhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	pricereviewparam "github.com/abolfazlnorzad/steve/param/pricereview"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetSetting(c echo.Context) error {
	resp, err := h.backofficeSvc.GetSetting(c.Request().Context(), pricereviewparam.GetSettingRequest{})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
