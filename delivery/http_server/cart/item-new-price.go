package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// SetItemUnitPriceBackoffice adjusts the effective unit price on one invoice
// line. The backend computes and stores the adjustment relative to the original
// price, supporting both discounts and markups.
func (h Handler) SetItemUnitPriceBackoffice(c echo.Context) error {
	const op = "carthandler.SetItemUnitPriceBackoffice"
	var req cartparam.SetItemUnitPriceBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateSetItemUnitPriceBackofficeRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.SetItemUnitPriceBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
