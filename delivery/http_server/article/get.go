package articlehandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) Get(c echo.Context) error {
	const op = "articlehandler.Get"
	id := c.Param("id")
	req := articleparam.GetArticleRequest{
		ID: id,
	}

	if err := h.vld.ValidateGetArticleRequest(req); err != nil {
		return err
	}

	resp, err := h.backofficeArticleSvc.GetArticle(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
