package orderhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) EditOrder(c echo.Context) error {
	const op = "orderhandler.EditOrder"

	var req orderparam.EditOrderRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(err)
	}
	req.OrderID = c.Param("id")
	req.EditedBy = claim.GetClaimsFromEchoContext(c).EntityID

	ctx := c.Request().Context()
	if err := h.vld.EditOrderValidator(ctx, req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.EditOrder(ctx, req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(h.translator.Translate(lang, "order.edited", nil), resp))
}
