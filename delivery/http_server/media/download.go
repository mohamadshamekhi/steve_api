package mediahandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetURL(c echo.Context) error {
	const op = "mediahandler.GetURL"
	id := c.Param("id")
	req := mediaparam.GetUrlRequest{Media: entity.Media{ID: id}}
	if err := h.vld.ValidateGetFile(req); err != nil {
		return err
	}

	res, err := h.svc.FindMediaByID(c.Request().Context(), mediaparam.GetMediaByIDRequest{ID: id})
	req.Media = res.Media
	resp, err := h.svc.GetURL(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
