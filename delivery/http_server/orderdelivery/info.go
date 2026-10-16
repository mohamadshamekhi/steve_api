package orderdeliveryhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderdeliveryparam "github.com/abolfazlnorzad/steve/param/orderdelivery"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetInfo(c echo.Context) error {
	resp, err := h.svc.GetInfo(c.Request().Context(), orderdeliveryparam.GetInfoRequest{Token: c.Param("token")})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(h.translator.Translate(lang, "global.action_successful", nil), resp))
}
