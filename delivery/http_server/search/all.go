package searchhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetSearchStorefrontWithProductsAndArticles(c echo.Context) error {
	const op = "searchhandler.searchhandler"
	search := c.QueryParam("search")

	req := searchparam.StorefrontSearchRequest{
		Search: commonparam.SearchRequest{Search: search},
	}

	resp, err := h.storefrontSearchSvc.SearchInProductAndBlog(c.Request().Context(), req)
	if err != nil {
		return err
	}

	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
