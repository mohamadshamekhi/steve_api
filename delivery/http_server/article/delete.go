package articlehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "articlehandler.Delete"
	id := c.Param("id")
	req := articleparam.DeleteArticleRequest{
		ID: id,
	}

	if err := h.vld.ValidateDeleteArticleRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeArticleSvc.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
