package backofficeauthhandler

import (
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	"github.com/abolfazlnorzad/steve/pkg/errmsg"
	httpmsg "github.com/abolfazlnorzad/steve/pkg/http_msg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) LoginByPhoneNumber(c echo.Context) error {
	var req backofficeauthparam.SendLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorMsgInvalidInput)
	}

	if fieldErrors, err := h.vld.ValidateLoginWithPhoneNumberRequest(req); err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	resp, err := h.svc.SendOTPCode(c.Request().Context(), req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, map[string]string{"message": msg})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) VerifyOTPCode(c echo.Context) error {
	var req backofficeauthparam.VerifyLoginCodeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return c.JSON(http.StatusBadRequest, errmsg.ErrorMsgInvalidInput)
	}
	if fieldErrors, err := h.vld.ValidateVerifyOTPCodeRequest(req); err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	resp, err := h.svc.VerifyOTPCode(c.Request().Context(), req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, map[string]string{"message": msg})
	}

	return c.JSON(http.StatusOK, resp)
}
