package producthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) UpdateAttrs(c echo.Context) error {
	const op = "productHandler.UpdateAttrs"
	var req productparam.AddProductAttributeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")

	if err := h.vld.ValidateUpdateSpecificationsRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.UpdateSpecifications(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) UpdateSingleAttrs(c echo.Context) error {
	const op = "producthandler.UpdateSingleAttrs"
	var req productparam.UpdateSingleAttributeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ProductID = c.Param("id")
	req.Key = c.Param("key")

	if err := h.vld.ValidateUpdateSingleSpecificationRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeProductSvc.UpdateSingleAttribute(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
