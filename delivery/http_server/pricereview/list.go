package pricereviewhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	pricereviewparam "github.com/abolfazlnorzad/steve/param/pricereview"
	"github.com/labstack/echo/v4"
)

func (h Handler) List(c echo.Context) error {
	var pagination commonparam.PaginationRequest
	_ = c.Bind(&pagination)

	req := pricereviewparam.ListNeedsReviewRequest{Pagination: pagination}
	resp, err := h.backofficeSvc.ListNeedsReview(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
