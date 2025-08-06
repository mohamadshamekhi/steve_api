package categoryhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	const op = "categoryhandler.Get"
	id := c.Param("id")
	req := categoryparam.GetCategoryRequest{
		ID: id,
	}

	if err := h.vld.ValidateGetCategoryRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeCategorySvc.GetCategory(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
