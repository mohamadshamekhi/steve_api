package orderdeliveryhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderdeliveryparam "github.com/abolfazlnorzad/steve/param/orderdelivery"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetInvoice(c echo.Context) error {
	resp, err := h.svc.GetInvoice(c.Request().Context(), orderdeliveryparam.GetInvoiceRequest{
		Token:        c.Param("token"),
		InvoiceToken: c.Request().Header.Get("X-Invoice-Token"),
	})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(h.translator.Translate(lang, "global.action_successful", nil), resp))
}
