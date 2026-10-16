package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/abolfazlnorzad/steve/delivery/http_server/response"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func NewCentralErrorHandler(translator *i18n.Translator) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		lang := getLanguageFromHeader(c.Request().Header)
		apiResponse, code := parseError(err, translator, lang)

		var richErr *richerror.RichError
		// ابتدا چک می‌کنیم آیا خطا از نوع ساختاریافته ماست؟
		if errors.As(err, &richErr) {
			logFields := []any{
				"error", err,
				"status", code,
				"request_id", c.Response().Header().Get(echo.HeaderXRequestID),
				"uri", c.Request().RequestURI,
				"method", c.Request().Method,
			}

			logger.L().Error("request failed with rich error",
				logFields...,
			)
		} else {
			// اگر نه، آن را به صورت یک خطای ساده لاگ می‌کنیم.
			logger.L().Error("request failed with simple error", "error", err)
		}

		if err := c.JSON(code, apiResponse); err != nil {
			logger.L().Error("failed to send json error response", "error", err)
		}
	}
}
func parseError(err error, translator *i18n.Translator, lang string) (response.APIResponse, int) {
	var richErr *richerror.RichError
	if errors.As(err, &richErr) {
		code := mapKindToHTTPStatusCode(richErr.Kind())
		mainMessage := translator.Translate(lang, richErr.MessageKey(), richErr.TemplateData())

		if summary, ok := richErr.GetMeta()["summary"]; ok {
			return response.NewErrorResponse(mainMessage, summary), code
		}

		if vErrs, ok := richErr.GetMeta()["fields"].(validation.Errors); ok {
			translatedFieldErrors := make(map[string]string)
			for field, fieldErr := range vErrs {
				var innerRichErr *richerror.RichError

				// الف) آیا خطای داخلی از نوع richerror است؟ (برای قوانین سفارشی پیچیده)
				if errors.As(fieldErr, &innerRichErr) {
					// پیام را با استفاده از کلید و داده‌های داینامیک خود آن ترجمه می‌کنیم.
					templateData := innerRichErr.GetMeta()
					if ids, dataOK := templateData["InvalidIDs"].([]string); dataOK {
						templateData["InvalidIDs"] = strings.Join(ids, ", ")
					}

					translatedFieldErrors[field] = translator.Translate(lang, innerRichErr.MessageKey(), templateData)

				} else if fieldErr != nil {
					// ب) ✅ اگر یک خطای ساده بود (مانند validation.Required)...
					// ۱. ابتدا نام خود فیلد را ترجمه می‌کنیم (مثلاً "PhoneNumber" -> "شماره تلفن")
					attributeKey := "attr." + field
					translatedAttribute := translator.Translate(lang, attributeKey, nil)

					// ۲. سپس پیام خطا را با نام ترجمه شده فیلد، می‌سازیم.
					templateData := map[string]any{
						"Attribute": translatedAttribute,
					}
					// پیام خطا (fieldErr.Error()) همان کلید i18n ماست (مثلاً "validation.required")
					finalMessage := translator.Translate(lang, fieldErr.Error(), templateData)

					translatedFieldErrors[field] = finalMessage
				}
			}

			return response.NewErrorResponse(mainMessage, translatedFieldErrors), code
		}

		//if fieldErrorKeys, ok := richErr.GetMeta()["fields"].(map[string]string); ok {
		//	translatedFieldErrors := make(map[string]string)
		//
		//	for field, key := range fieldErrorKeys {
		//		attributeKey := "attr." + field
		//		translatedAttribute := translator.Translate(lang, attributeKey, nil)
		//
		//		templateData := map[string]any{
		//			"Attribute": translatedAttribute,
		//		}
		//		finalMessage := translator.Translate(lang, key, templateData)
		//
		//		translatedFieldErrors[field] = finalMessage
		//	}
		//
		//	return response.NewErrorResponse(mainMessage, translatedFieldErrors), code
		//}

		if code >= http.StatusInternalServerError {
			mainMessage = translator.Translate(lang, "global.something_wrong", nil)
			return response.NewErrorResponse(mainMessage, nil), code
		}

		return response.NewErrorResponse(mainMessage, nil), code
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		msg := translator.Translate(lang, err.Error(), nil)
		return response.NewErrorResponse(msg, nil), echoErr.Code
	}

	msg := translator.Translate(lang, "global.something_wrong", nil)
	return response.NewErrorResponse(msg, nil), http.StatusInternalServerError
}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindBadRequest:
		return http.StatusBadRequest
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindInvalid:
		return http.StatusBadRequest
	case richerror.KindUnauthorized:
		return http.StatusUnauthorized
	case richerror.KindConflict:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func getLanguageFromHeader(header http.Header) string {
	if header.Get("Accept-Language") == "en" {
		return "en"
	}
	return "fa"
}
