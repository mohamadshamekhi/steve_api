package producthandler

import (
	"fmt"
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
	fmt.Println("op", op)
	resp, err := h.storefrontProductSvc.GetProductsWithAutoDiscount(c.Request().Context(), productparam.GetProductsWithAutoDiscountsRequest{})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h Handler) GetProductsWithSpecialDiscount(c echo.Context) error {
	const op = "productservice.GetProductsWithSpecialDiscount"
	fmt.Println("op", op)
	resp, err := h.storefrontProductSvc.GetProductsWithSpecialDiscount(c.Request().Context(), productparam.GetProductsWithSpecialDiscountRequest{})
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
	req := productparam.StorefrontGetAllProductRequest{
		Filter: productparam.StorefrontFilterProduct{
			Filter: make(map[string]any),
		},
	}

	queryParams := c.QueryParams()

	if v := queryParams.Get("brand_id"); v != "" {
		req.Filter.BrandID = &v
	}
	if v := queryParams.Get("category_id"); v != "" {
		req.Filter.CategoryID = &v
	}
	if v := queryParams.Get("category_slug"); v != "" {
		req.Filter.CategorySlug = &v
	}
	if v := queryParams.Get("color"); v != "" {
		req.Filter.Color = &v
	}
	if v := queryParams.Get("order_by"); v != "" {
		req.Filter.OrderBy = &v
	}

	if v := queryParams.Get("is_second_hand"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			req.Filter.IsSecondHand = &b
		}
	}
	if v := queryParams.Get("is_exist"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			req.Filter.IsExist = &b
		}
	}
	if v := queryParams.Get("min_price"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			req.Filter.MinPrice = &p
		}
	}
	if v := queryParams.Get("max_price"); v != "" {
		if p, err := strconv.ParseUint(v, 10, 64); err == nil {
			req.Filter.MaxPrice = &p
		}
	}

	// ۳. مدیریت صفحه‌بندی
	pageNumber, _ := strconv.ParseUint(queryParams.Get("page_number"), 10, 64)
	if pageNumber == 0 {
		pageNumber = commonparam.DefaultPageNumber
	}
	pageSize, _ := strconv.ParseUint(queryParams.Get("page_size"), 10, 64)
	if pageSize == 0 {
		pageSize = commonparam.DefaultPageSize
	}
	req.Pagination.PageSize = uint(pageSize)
	req.Pagination.PageNumber = uint(pageNumber)

	reservedKeys := map[string]bool{
		"brand_id":       true,
		"category_id":    true,
		"category_slug":  true,
		"color":          true,
		"order_by":       true,
		"is_second_hand": true,
		"is_exist":       true,
		"min_price":      true,
		"max_price":      true,
		"page_number":    true,
		"page_size":      true,
	}

	for key, values := range queryParams {
		if !reservedKeys[key] && len(values) > 0 {
			req.Filter.Filter[key] = values[0]
		}
	}

	resp, err := h.storefrontProductSvc.GetFilteredProduct(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
