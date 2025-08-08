package brandhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "brandhandler.Delete"
	id := c.Param("id")
	req := brandparam.DeleteBrandRequest{
		ID: id,
	}

	if err := h.validator.ValidateDeleteBrandRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
