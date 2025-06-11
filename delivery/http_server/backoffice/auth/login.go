package backofficeauthhandler

import (
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) LoginByPhoneNumber(c echo.Context) error {
	const op = "backofficeauthhandler.LoginByPhoneNumber"
	var req backofficeauthparam.SendLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	if err := h.vld.ValidateLoginWithPhoneNumberRequest(req); err != nil {
		return err
	}
	resp, err := h.svc.SendOTPCode(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) VerifyOTPCode(c echo.Context) error {
	const op = "backofficeauthhandler.VerifyOTPCode"
	var req backofficeauthparam.VerifyLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	if err := h.vld.ValidateVerifyOTPCodeRequest(req); err != nil {
		return err
	}
	resp, err := h.svc.VerifyOTPCode(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
