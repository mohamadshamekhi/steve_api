package userhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	userparam "github.com/abolfazlnorzad/steve/param/user"
	"github.com/labstack/echo/v4"
)

func (h Handler) Me(c echo.Context) error {
	var req userparam.MeRequest
	req.UserID = c.Get("user_id").(string)
	resp, err := h.storeFrontSvc.Me(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)

	successMessage := h.translator.Translate(lang, "global.action_successful", nil)

	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
