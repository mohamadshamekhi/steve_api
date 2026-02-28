package orderhandler

import (
	"context"
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Register(c echo.Context) error {
	const op = "orderhandler.Register"
	var req orderparam.RegisterOrderRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	cl := claim.GetClaimsFromEchoContext(c)
	c.Set("user_id", cl.EntityID)

	if err := h.vld.RegisterOrderValidator(c.Request().Context(), req); err != nil {
		return err
	}
	ctx := context.WithValue(c.Request().Context(), "user_id", cl.EntityID)
	resp, err := h.storefrontSvc.RegisterOrder(ctx, req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) RegisterBackoffice(c echo.Context) error {
	const op = "orderhandler.RegisterBackoffice"
	var req orderparam.RegisterOrderBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.RegisterOrderBackofficeValidator(c.Request().Context(), req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.RegisterOrderBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) CompleteOrderBackoffice(c echo.Context) error {
	const op = "orderhandler.CompleteOrderBackoffice"
	var req orderparam.CompletedOrderBackofficeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	resp, err := h.backofficeSvc.CompleteOrderBackoffice(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
