package userhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
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

func (h Handler) UpdateUserBackoffice(c echo.Context) error {
	const op = "userhandler.UpdateUserBackoffice"
	var req userparam.UpdateUserBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	req.UserID = c.Param("id")

	if err := h.vld.ValidateBackofficeUpdateUserRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.UpdateUser(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
