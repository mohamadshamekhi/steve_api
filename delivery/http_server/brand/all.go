package brandhandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "brandhandler.All"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	name := c.QueryParam("name")

	req := brandparam.GetAllBrandsParamRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: []map[string]any{{"name": name}},
		},
	}

	isFavoriteStr := c.QueryParam("is_favorite")
	if isFavoriteStr != "" {
		isFavorite, bErr := strconv.ParseBool(isFavoriteStr)
		if bErr != nil {
			return richerror.New(op).
				WithKind(richerror.KindInvalid).
				WithMessageKey("validation.invalid_input").
				WithWarpError(bErr)
		}
		req.Filter.Filter = append(req.Filter.Filter, map[string]any{"is_favorite": isFavorite})
	}

	resp, err := h.backofficeSvc.All(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
func (h Handler) GetAllStorefront(c echo.Context) error {
	const op = "brandhandler.GetAllStorefront"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}

	name := c.QueryParam("name")
	isFavorit := c.QueryParam("is_favorit")

	filters := []map[string]any{}
	if name != "" {
		filters = append(filters, map[string]any{
			"name": name,
		})
	}
	if isFavorit != "" {
		fav, err := strconv.ParseBool(isFavorit)
		fmt.Println("fav", fav)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid is_favorit value")
		}
		filters = append(filters, map[string]any{
			"is_favorit": fav,
		})
	}

	req := brandparam.GetAllStorefrontBrandsParamRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
		Filter: commonparam.FilterRequest{
			Filter: filters,
		},
	}

	resp, err := h.storefrontSvc.GetAllBrands(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
