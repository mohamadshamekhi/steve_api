package producthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	const op = "producthandler.Get"
	id := c.Param("id")
	req := productparam.GetProductRequest{
		ID: id,
	}

	if err := h.vld.ValidateGetProductRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeProductSvc.GetProduct(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) StorefrontGetProduct(c echo.Context) error {
	const op = "producthandler.StorefrontGetProduct"
	id := c.Param("id")
	req := productparam.StorefrontGetProductRequest{
		ID: id,
	}

	if err := h.vld.ValidateStorefrontGetProductRequest(req); err != nil {
		return err
	}

	resp, err := h.storefrontProductSvc.GetProduct(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
func (h Handler) StorefrontGetProductBySlug(c echo.Context) error {
	const op = "producthandler.StorefrontGetProductBySlug"
	slug := c.Param("slug")
	req := productparam.StorefrontGetProductBySlugRequest{
		Slug: slug,
	}

	if err := h.vld.ValidateStorefrontGetProductBySlugRequest(req); err != nil {
		return err
	}

	resp, err := h.storefrontProductSvc.GetProductBySlug(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) GetMobileList(c echo.Context) error {
	const op = "producthandler.GetMobileList"
	name := c.QueryParam("name")
	brandID := c.QueryParam("brand_id")
	req := productparam.GetFilterMobileListWithMainInfoRequest{}
	if name != "" {
		req.Name = &name
	}
	if brandID != "" {
		req.BrandID = &brandID
	}

	resp, err := h.storefrontProductSvc.GetFilterMobileListWithMainInfo(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
