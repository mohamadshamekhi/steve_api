package commenthandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) UpdateComment(c echo.Context) error {
	const op = "commentHandler.UpdateComment"
	var req commentparam.UpdateCommentRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.CommentID = c.Param("id")
	if err := h.validator.ValidateUpdateCommentRequest(req); err != nil {
		return err
	}

	resp, err := h.svc.UpdateComment(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
