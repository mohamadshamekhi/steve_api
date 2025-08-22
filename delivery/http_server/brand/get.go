package brandhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	const op = "brandhandler.Get"
	id := c.Param("id")
	req := brandparam.GetBrandRequest{
		ID: id,
	}

	if err := h.validator.ValidateGetBrandRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.GetBrand(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
