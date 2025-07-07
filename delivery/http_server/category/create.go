package categoryhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	categoryparam "github.com/abolfazlnorzad/steve/param/category"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Create(c echo.Context) error {
	const op = "categoryhandler.Create"
	var req categoryparam.CreateCategoryRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateCreateCategoryRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeCategorySvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
