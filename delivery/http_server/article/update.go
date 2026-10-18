package articlehandler

import (
	"net/http"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	"github.com/abolfazlnorzad/steve/pkg/claim"
	"github.com/abolfazlnorzad/steve/pkg/key"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/labstack/echo/v4"
)

func (h Handler) Update(c echo.Context) error {
	const op = "articlehandler.Update"

	var req articleparam.UpdateArticleRequest
	if bErr := c.Bind(&req); bErr != nil {
		return richerror.New(op).
			WithKind(richerror.KindInvalid).
			WithMessageKey("validation.invalid_input").
			WithWarpError(bErr)
	}
	req.ID = c.Param("id")
	if err := h.vld.ValidateUpdateArticleRequest(req); err != nil {
		return err
	}

	// Publishing is a distinct, deliberate action (key.PublishArticle) even
	// though it rides on the generic update endpoint: a caller that only has
	// article:update (e.g. the SEO/content agent) may edit everything else
	// about a draft but can never flip is_published itself.
	if req.IsPublished != nil {
		claims := claim.GetClaimsFromEchoContext(c)
		allowed, aErr := h.aclSvc.CheckAccess(c.Request().Context(), claims.EntityID, key.PublishArticle)
		if aErr != nil {
			return richerror.New(op).WithWarpError(aErr)
		}
		if !allowed {
			return richerror.New(op).WithKind(richerror.KindForbidden).WithMessageKey("global.forbidden")
		}
	}

	resp, err := h.backofficeArticleSvc.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}
	lang := getLanguageFromHeader(c.Request().Header)
	successMessage := h.translator.Translate(lang, "global.action_successful", nil)
	apiResponse := response.NewSuccessResponse(successMessage, resp)

	return c.JSON(http.StatusOK, apiResponse)
}
