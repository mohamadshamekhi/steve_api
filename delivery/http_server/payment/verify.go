package paymenthandler

import (
	"context"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	paymentparam "github.com/abolfazlnorzad/steve/param/payment"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Verify(c echo.Context) error {
	const op = "paymentHandler.Verify"
	var req paymentparam.VerifyPaymentRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	cl := claim.GetClaimsFromEchoContext(c)
	c.Set("user_id", cl.EntityID)

	ctx := context.WithValue(c.Request().Context(), "user_id", cl.EntityID)
	resp, err := h.svc.VerifyPayment(ctx, req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
