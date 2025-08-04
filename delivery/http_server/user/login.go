package userhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userauthparam "github.com/abolfazlnorzad/steve/param/user/auth"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Login(c echo.Context) error {
	const op = " userhandler.Login"
	var req userauthparam.UserSendLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	if err := h.vld.ValidateUserLoginWithPhoneNumberRequest(req); err != nil {
		return err
	}
	resp, err := h.storeFrontSvc.SendOTPCode(c.Request().Context(), req)
	if err != nil {
		return err
	}

	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "auth.otp_send", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) VerifyOTPCode(c echo.Context) error {
	const op = "userhandler.VerifyOTPCode"
	var req userauthparam.UserVerifyLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	if err := h.vld.ValidateUserVerifyOTPCodeRequest(req); err != nil {
		return err
	}
	resp, err := h.storeFrontSvc.VerifyOTPCode(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "auth.login_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
