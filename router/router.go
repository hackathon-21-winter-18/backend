package router

import (
	"net/http"

	sess "github.com/hackathon-21-winter-18/backend/session" // sessだけ使うって意味？
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4/middleware"
)

var s sess.Session

func SetRouting(sess sess.Session) {
	s = sess

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
			apiOauth.POST("/logout", postLogout, userAuthMiddleware)
			apiOauth.GET("/whoamI", getWhoamI, userAuthMiddleware)
		}

		apiPalaces := api.Group("/palaces")
		{
			apiPalaces.GET("/me/:userID", getPalaces, userAuthMiddleware)
			apiPalaces.POST("/me/:userID", postPalace, userAuthMiddleware)
			apiPalaces.PUT("/:palaceID", putPalace, userAuthMiddleware)
			apiPalaces.DELETE("/:palaceID", deletePalace, userAuthMiddleware)
		}


	}

	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
