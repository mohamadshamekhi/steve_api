package userhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) UpdateProfile(c echo.Context) error {
	const op = "userhandler.UpdateProfile"
	var req userparam.StoreFrontUpdateProfileRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	if err := h.vld.ValidateStoreFrontUpdateProfileRequest(req); err != nil {
		return err
	}
	resp, err := h.storeFrontSvc.UpdateProfile(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
