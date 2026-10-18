package backofficehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	backofficeparam "github.com/abolfazlnorzad/steve/param/backoffice"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) IssueServiceToken(c echo.Context) error {
	const op = "backofficehandler.IssueServiceToken"
	id := c.Param("id")
	if id == "" {
		return richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("validation.required")
	}
	req := backofficeparam.IssueServiceTokenRequest{ID: id}

	resp, err := h.svc.IssueServiceToken(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
