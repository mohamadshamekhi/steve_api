package shippinghandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Create(c echo.Context) error {
	const op = "settinghandler.Create"

	var req settingparam.AddShippingRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateAddShippingRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
