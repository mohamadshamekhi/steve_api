package pricereviewhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	pricereviewparam "github.com/abolfazlnorzad/steve/param/pricereview"
	"github.com/labstack/echo/v4"
)

func (h Handler) Confirm(c echo.Context) error {
	req := pricereviewparam.ConfirmRequest{ProductID: c.Param("id")}

	resp, err := h.backofficeSvc.Confirm(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
