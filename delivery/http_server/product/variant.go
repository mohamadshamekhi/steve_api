package producthandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) AddVariant(c echo.Context) error {
	const op = "productHandler.AddVariant"
	var req productparam.AddVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateAddVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.AddVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) DeleteVariant(c echo.Context) error {
	const op = "productHandler.DeleteVariant"
	var req productparam.DeleteVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateDeleteVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.DeleteVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) UpdateVariant(c echo.Context) error {
	const op = "productHandler.UpdateVariant"
	var req productparam.UpdateVariantRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateUpdateVariantRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.UpdateVariant(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
