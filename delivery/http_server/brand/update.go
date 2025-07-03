package brandhandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Update(c echo.Context) error {
	const op = "brandHandler.Update"

	var req brandparam.UpdateBrandRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.validator.ValidateUpdateBrandRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeSvc.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
