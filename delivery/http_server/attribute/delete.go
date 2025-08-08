package attributehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "attributehandler.Delete"
	id := c.Param("id")
	req := attributeparam.DeleteAttributeRequest{
		ID: id,
	}

	if err := h.vld.ValidateDeleteAttributeRequest(req); err != nil {
		return err
	}

	resp, err := h.svc.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
