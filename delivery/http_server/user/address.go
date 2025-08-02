package userhandler

import (
	"fmt"
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) AddAddress(c echo.Context) error {
	const op = "userhandler.AddAddress"
	var req userparam.StoreFrontAddAddressRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	req.UserID = cl.EntityID
	fmt.Println("req: ", req)
	if err := h.vld.ValidateStoreFrontAddAddressRequest(req); err != nil {
		return err
	}
	resp, err := h.storeFrontSvc.AddAddress(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) UpdateAddress(c echo.Context) error {
	const op = "userhandler.UpdateAddress"
	var req userparam.StoreFrontUpdateAddressRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(bErr)
	}
	cl := claim.GetClaimsFromEchoContext(c)
	req.UserID = cl.EntityID
	req.ID = c.Param("id")
	if err := h.vld.ValidateStoreFrontUpdateAddressRequest(req); err != nil {
		return err
	}
	resp, err := h.storeFrontSvc.UpdateAddress(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
