package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// SetItemManualDiscountBackoffice records the discount given by hand on one
// line of an in-person invoice.
func (h Handler) SetItemManualDiscountBackoffice(c echo.Context) error {
	const op = "carthandler.SetItemManualDiscountBackoffice"
	var req cartparam.SetItemManualDiscountBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateSetItemManualDiscountBackofficeRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.SetItemManualDiscountBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}

// SetCartManualDiscountBackoffice records the discount given by hand on the
// in-person invoice as a whole.
func (h Handler) SetCartManualDiscountBackoffice(c echo.Context) error {
	const op = "carthandler.SetCartManualDiscountBackoffice"
	var req cartparam.SetCartManualDiscountBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateSetCartManualDiscountBackofficeRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.SetCartManualDiscountBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
