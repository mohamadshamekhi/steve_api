package attributedefinitionhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "attributedefinitionhandler.Delete"
	id := c.Param("id")
	req := attributedefinitionparam.DeleteAttributeDefinitionRequest{
		ID: id,
	}

	if err := h.vld.ValidateDeleteAttributeDefinitionRequest(req); err != nil {
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
