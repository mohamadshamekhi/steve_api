package carthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// UpsertCustomItemBackoffice adds or edits a hand-typed line on an in-person
// invoice — a repair fee, an Apple ID, or a part bought from another shop.
func (h Handler) UpsertCustomItemBackoffice(c echo.Context) error {
	const op = "carthandler.UpsertCustomItemBackoffice"
	var req cartparam.UpsertCustomItemBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateUpsertCustomItemBackofficeRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.UpsertCustomItemBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
