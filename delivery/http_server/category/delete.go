package categoryhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "categoryhandler.Delete"
	id := c.Param("id")
	req := categoryparam.DeleteCategoryRequest{
		ID: id,
	}

	if err := h.vld.ValidateDeleteCategoryRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeCategorySvc.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
