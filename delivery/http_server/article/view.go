package articlehandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

// visitorCookieName holds an anonymous, per-browser identifier used to
// dedupe article view counts. It carries no personal data - just a random
// ID minted on first visit.
const visitorCookieName = "steve_visitor_id"

// visitorID returns the caller's existing visitor cookie, or mints and sets
// a new one if none is present yet.
func (h Handler) visitorID(c echo.Context) string {
	if cookie, err := c.Cookie(visitorCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	id := ulid.Make().String()
	c.SetCookie(&http.Cookie{
		Name:     visitorCookieName,
		Value:    id,
		Path:     "/",
		Expires:  time.Now().AddDate(2, 0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return id
}
