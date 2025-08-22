package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetSingle(c echo.Context) error {
	const op = "carthandler.GetSingle"
	id := c.Param("id")
	req := cartparam.GetCartRequest{
		ID: id,
	}

	if err := h.vld.ValidateGetCartRequest(req); err != nil {
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
