package bannerhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	const op = "bannerhandler.Delete"

	var req settingparam.DeleteBannerRequest
	id := c.Param("id")
	req.ID = id

	resp, err := h.backofficeSvc.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
