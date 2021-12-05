package router

import (
	"net/http"

	sess "github.com/hackathon-21-winter-18/backend/session" // sessだけ使うって意味？
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
			apiOauth.POST("/signup", postSignUp)
			apiOauth.POST("/login", postLogin)
			apiOauth.GET("/whoamI", getWhoamI)
		}

		apiPalaces := api.Group("/palaces")
		{
			apiPalaces.POST("/:userID", postPalace)
		}


	}

	err := e.Start(":3000")
	if err != nil {
		panic(err)
	}
}
