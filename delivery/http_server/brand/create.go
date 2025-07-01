package brandhandler

import (
	"errors"
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	brandparam "github.com/abolfazlnorzad/steve/param/brand"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Create(c echo.Context) error {
	const op = "brandhandler.Create"
	var req brandparam.CreateBrandRequest
	if bErr := c.Bind(&req); bErr != nil {
		var echoErr *echo.HTTPError
		if errors.As(bErr, &echoErr) {
			if echoErr.Code == http.StatusUnsupportedMediaType {
				return richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("media.unsupported_media_type").
					WithWarpError(bErr)
			}
		}
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}

	if err := h.validator.ValidateCreateBrandRequest(req); err != nil {
		return err
	}
	resp, err := h.backofficeSvc.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
