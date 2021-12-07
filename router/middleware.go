package router

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func userAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			// TODO
			return errSessionNotFound(err)
		}

		userID := sess.Values["userID"]
		if sess.Values["userID"] == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		c.Set("userID", userID)

		return next(c)
	}
}
