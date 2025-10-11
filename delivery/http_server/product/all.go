package producthandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "producthandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	ENTitle := c.QueryParam("en_title")
	name := c.QueryParam("name")
	isPub := c.QueryParam("is_published")
	isSecondHand := c.QueryParam("is_second_hand")
	brandID := c.QueryParam("brand_id")
	categoryID := c.QueryParam("category_id")
	req := productparam.GetAllProductsRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"en_title": ENTitle}, {"name": name}, {"brand_id": brandID},
				{"category_id": categoryID}, {"is_published": isPub}, {"is_second_hand": isSecondHand}},
		},
	}

	resp, err := h.backofficeProductSvc.GetAllProducts(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}

func (h Handler) StorefrontGetProductsWithAutoDiscount(c echo.Context) error {
	const op = "productservice.StorefrontGetProductsWithAutoDiscount"
	resp, err := h.storefrontProductSvc.GetProductsWithAutoDiscount(c.Request().Context(), productparam.GetProductsWithAutoDiscountsRequest{})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) StorefrontGetFilterableProduct(c echo.Context) error {
	const op = "productservice.StorefrontGetFilterableProduct"
	req := productparam.StorefrontGetAllProductRequest{}

	brandID := c.QueryParam("brand_id")
	if brandID != "" {
		req.Filter.BrandID = &brandID
	}

	categoryID := c.QueryParam("category_id")
	if categoryID != "" {
		req.Filter.CategoryID = &categoryID
	}

	color := c.QueryParam("color")
	if color != "" {
		req.Filter.Color = &color
	}

	orderBy := c.QueryParam("order_by")
	if orderBy != "" {
		req.Filter.OrderBy = &orderBy
	}

	if minPriceStr := c.QueryParam("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseUint(minPriceStr, 10, 64); err == nil {
			req.Filter.MinPrice = &minPrice
		} else {
			return err
		}
	}

	if maxPriceStr := c.QueryParam("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseUint(maxPriceStr, 10, 64); err == nil {
			req.Filter.MaxPrice = &maxPrice
		} else {
			return err
		}
	}

	if isExistStr := c.QueryParam("is_exist"); isExistStr != "" {
		if isExist, err := strconv.ParseBool(isExistStr); err == nil {
			req.Filter.IsExist = &isExist
		} else {

			return err
		}
	}
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}

	req.Pagination.PageSize = uint(pageSize)
	req.Pagination.PageNumber = uint(pageNumber)

	resp, err := h.storefrontProductSvc.GetFilteredProduct(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
