package bannerhandler

import (
	"net/http"
	"strconv"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	settingparam "github.com/abolfazlnorzad/steve/param/setting"
	"github.com/labstack/echo/v4"
)

func (h Handler) All(c echo.Context) error {
	const op = "bannerhandler.All"

	isActive := c.QueryParam("is_active")
	var isActiveBool *bool
	if isActive != "" {
		b, _ := strconv.ParseBool(isActive)
		isActiveBool = &b

	} else {
		isActiveBool = nil
	}
	req := settingparam.ListBannerRequest{
		IsActive: isActiveBool,
	}

	resp, err := h.backofficeSvc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)

}
