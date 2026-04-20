package orderhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	orderparam "github.com/abolfazlnorzad/steve/param/order"
	"github.com/labstack/echo/v4"
)

func (h Handler) CheckPendingExist(c echo.Context) error {
	const op = "orderhandler.CheckPendingExist"

	var req orderparam.CheckPendingExistRequest

	req.UserID = c.Param("id")

	resp, err := h.backofficeSvc.CheckPendingExist(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
