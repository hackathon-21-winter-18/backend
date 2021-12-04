package router

import (
	"net/http"

	sess "github.com/hackathon-winter-18/backend/session" // sessだけ使うって意味？
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4/middleware"
)

func SetRouting(sess sess.Session) {
	// s := sess

	e := echo.New()
	e.Use(session.Middleware(sess.Store()))
	e.Use(middleware.Logger())

	api := e.Group("/api")
	{
		apiPing := api.Group("/ping")
		{
			apiPing.GET("", func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusOK, "pong!")
			})
		}

		apiOauth := api.Group("/oauth")
		{
			apiOauth.POST("/singup", postSignUp)
			apiOauth.POST("/login", postLogin)
		}
	}

	err := e.Start(":3000")
	if err != nil {
		panic(err)
	}
}
