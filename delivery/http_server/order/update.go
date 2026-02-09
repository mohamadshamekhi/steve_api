package orderhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) UpdateStatus(c echo.Context) error {
	const op = "orderhandler.UpdateStatus"

	var req orderparam.UpdateStatusOrderRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(err)
	}

	id := c.Param("id")
	req.OrderID = id

	ctx := c.Request().Context()
	if err := h.vld.UpdateStatusOrderValidator(ctx, req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.UpdateStatus(ctx, req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
