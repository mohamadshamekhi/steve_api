package attributedefinitionhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	attributedefinitionparam "github.com/abolfazlnorzad/steve/param/attribute_definition"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Update(c echo.Context) error {
	const op = "attributedefinitionhandler.Update"

	var req attributedefinitionparam.UpdateAttributeDefinitionRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ID = c.Param("id")

	if err := h.vld.ValidateUpdateAttributeDefinitionRequest(req); err != nil {
		return err
	}

	resp, err := h.svc.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
