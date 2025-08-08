package backofficeauthhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
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

	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "auth.otp_send", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
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
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "auth.login_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
