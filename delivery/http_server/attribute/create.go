package attributehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	attributeparam "github.com/abolfazlnorzad/steve/param/attribute"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Create(c echo.Context) error {
	const op = "attributehandler.Create"
	var req attributeparam.AddAttributeRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.vld.ValidateCreateAttributeRequest(req); err != nil {
		return err
	}
	resp, err := h.svc.CreateAttribute(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
