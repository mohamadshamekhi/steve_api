package aclhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	aclparam "github.com/abolfazlnorzad/steve/param/acl"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetAllRole(c echo.Context) error {
	const op = "aclhandler.GetAllRole"
	search := c.QueryParam("search")

	resp, err := h.roleSvc.GetAllRoles(c.Request().Context(), aclparam.GetAllRolesRequest{
		Search: search,
	})
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
