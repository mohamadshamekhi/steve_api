package commenthandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetComments(c echo.Context) error {
	const op = "commentHandler.GetComments"
	pageNumber, nErr := strconv.ParseUint(c.QueryParam("page_number"), 10, 64)
	if nErr != nil {
		pageNumber = commonparam.DefaultPageNumber
	}

	pageSize, sErr := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if sErr != nil {
		pageSize = commonparam.DefaultPageSize
	}
	req := commentparam.StorefrontGetCommentsRequest{
		Pagination: commonparam.PaginationRequest{
			PageSize:   uint(pageSize),
			PageNumber: uint(pageNumber),
		},
	}
	req.CommentableID = c.Param("commentableid")
	req.CommentableType = c.Param("commentabletype")

	resp, err := h.storefrontSvc.GetAprovedCommentsWithChildren(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
