package aclhandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	aclparam "github.com/abolfazlnorzad/steve/param/acl"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetAllPermissions(c echo.Context) error {
	const op = "aclhandler.GetAllPermissions"

	search := c.QueryParam("search")

	resp, err := h.permissionSvc.GetAllPermissions(c.Request().Context(), aclparam.GetAllPermissionsRequest{
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
