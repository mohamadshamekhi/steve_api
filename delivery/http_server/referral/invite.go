package referralhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	referralparam "github.com/abolfazlnorzad/steve/param/referral"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

// InviteByPhone records an outstanding invite for a phone number and texts the
// invite link to it.
func (h Handler) InviteByPhone(c echo.Context) error {
	const op = "referralhandler.InviteByPhone"
	var req referralparam.InviteByPhoneRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateInviteByPhoneRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.svc.InviteByPhone(c.Request().Context(), cl.EntityID, req)
	if err != nil {
		return err
	}

	lang := getLanguageFromHeader(c.Request().Header)
	msgKey := "referral.invite_sent"
	if !resp.SMSSent {
		msgKey = "referral.invite_recorded_sms_failed"
	}
	msg := h.translator.Translate(lang, msgKey, nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, resp))
}
