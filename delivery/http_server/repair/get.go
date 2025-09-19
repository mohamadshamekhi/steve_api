package repairhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	repairparam "github.com/abolfazlnorzad/steve/param/repair"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(c echo.Context) error {
	id := c.Param("id")
	req := repairparam.GetRepairRequest{
		ID: id,
	}

	resp, err := h.backofficeSvc.GetRepair(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
