package carthandler

import (
	"context"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	cartparam "github.com/abolfazlnorzad/steve/param/cart"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) DeleteItem(c echo.Context) error {
	const op = "handlercart.DeleteItem"

	claims := claim.GetClaimsFromEchoContext(c)

	userID := claims.EntityID

	var req cartparam.DeleteItemFromCartRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(err)
	}
	req.UserID = userID
	ctx := context.WithValue(c.Request().Context(), "user_id", userID)
	if err := h.vld.ValidateDeleteItemRequest(req); err != nil {
		return err
	}
	resp, err := h.storefrontSvc.DeleteItemFromCart(ctx, req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) DeleteItemBackoffice(c echo.Context) error {
	const op = "handlercart.DeleteItemBackoffice"
	var req cartparam.DeleteItemFromCartBackofficeRequest
	if err := c.Bind(&req); err != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("global.invalid_input").
			WithWarpError(err)
	}
	if err := h.vld.ValidateDeleteItemBackofficeRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.DeleteItemFromCartBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
