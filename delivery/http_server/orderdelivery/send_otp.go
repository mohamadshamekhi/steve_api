package orderdeliveryhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderdeliveryparam "github.com/abolfazlnorzad/steve/param/orderdelivery"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) SendOTP(c echo.Context) error {
	const op = "orderdeliveryhandler.SendOTP"

	var req orderdeliveryparam.SendOTPRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.Token = c.Param("token")

	if err := h.vld.SendOTP(req); err != nil {
		return err
	}
	resp, err := h.svc.SendOTP(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(h.translator.Translate(lang, "order.delivery_otp_sent", nil), resp))
}
