package backofficehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) AddNewBackofficeUser(c echo.Context) error {
	const op = "backofficehandler.AddNewBackofficeUser"
	var req backofficeparam.AddNewBackofficeUserRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateAddNewBackofficeUserRequest(req); err != nil {
		return err
	}
	resp, err := h.svc.AddNewBackofficeUser(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
