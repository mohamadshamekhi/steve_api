package mediahandler

import (
	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) UploadFile(c echo.Context) error {
	const op = "mediahandler.UploadFile"

	header, err := c.FormFile("file")
	if err != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("media.file_not_provided").
			WithWarpError(err)
	}

	file, err := header.Open()
	if err != nil {
		return richerror.New(op).WithKind(richerror.KindUnexpected).WithWarpError(err).WithMessageKey("global.something_wrong")
	}
	defer file.Close()

	var req mediaparam.UploadMediaRequest
	req.Domain = c.FormValue("domain")
	req.UploaderID = c.FormValue("uploader_id")
	req.UploaderType = c.FormValue("uploader_type")
	req.IsPublic = c.FormValue("is_public") == "true"
	req.UploaderType = "backoffice"
	cl := claim.GetClaimsFromEchoContext(c)
	req.UploaderID = cl.EntityID

	req.File = file
	req.Header = header
	if err := h.vld.ValidateUploadRequest(req); err != nil {
		return err
	}

	resp, err := h.svc.UploadFile(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
