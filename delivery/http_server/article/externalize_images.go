package articlehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/labstack/echo/v4"
)

// ExternalizeInlineImages is a one-off maintenance action: it uploads the
// base64 images still embedded in existing article bodies.
func (h Handler) ExternalizeInlineImages(c echo.Context) error {
	resp, err := h.backofficeArticleSvc.ExternalizeInlineImages(c.Request().Context(), articleparam.ExternalizeInlineImagesRequest{})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	return c.JSON(http.StatusOK, response.NewSuccessResponse(successMessage, resp))
}
