package sellhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	sellparam "github.com/abolfazlnorzad/steve/param/sell"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	id := c.Param("id")
	req := sellparam.GetSellRequest{
		ID: id,
	}

	resp, err := h.backofficeSvc.GetSell(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
