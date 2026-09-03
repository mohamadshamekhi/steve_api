package partnerhandler

import (
	"context"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	partnerparam "github.com/abolfazlnorzad/steve/param/partner"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Apply(c echo.Context) error {
	const op = "partnerhandler.Apply"
	var req partnerparam.ApplyRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").WithWarpError(bErr)
	}
	if err := h.validator.ValidateApplyRequest(req); err != nil {
		return err
	}

	cl := claim.GetClaimsFromEchoContext(c)
	ctx := context.WithValue(c.Request().Context(), "user_id", cl.EntityID)

	resp, err := h.storefrontSvc.Apply(ctx, cl.EntityID, req)
	if err != nil {
		return err
	}

	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "partner.application_submitted", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, resp))
}

func (h Handler) Me(c echo.Context) error {
	cl := claim.GetClaimsFromEchoContext(c)
	resp, err := h.storefrontSvc.Me(c.Request().Context(), cl.EntityID)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	msg := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(msg, resp))
}
