package discounthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	discountparam "github.com/abolfazlnorzad/steve/param/discount"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	const op = "discounthandler.Get"
	id := c.Param("id")
	req := discountparam.GetSingleDiscountRequest{
		ID: id,
	}

	if err := h.validator.ValidateGetDiscountRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.GetSingle(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
