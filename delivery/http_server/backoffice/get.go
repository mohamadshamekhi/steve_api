package backofficehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	const op = "backofficehandler.Get"
	id := c.Param("id")
	req := backofficeparam.GetBackofficeRequest{
		ID: id,
	}
	if err := h.vld.ValidateGetBackofficeRequest(req); err != nil {
		return err
	}

	resp, err := h.svc.GetBackoffice(c.Request().Context(), req)
	if err != nil {
		return err

	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
