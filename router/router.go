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
	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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
			apiOauth.POST("/logout", postLogout)
			apiOauth.GET("/whoamI", getWhoamI)
		}

		apiPalaces := api.Group("/palaces")
		{
			apiPalaces.GET("/palaces", getPalaces)
			apiPalaces.GET("/me/:userID", getMyPalaces)
			apiPalaces.POST("/me/:userID", postPalace)
			apiPalaces.PUT("/:palaceID", putPalace, userAuthMiddleware)
			apiPalaces.DELETE("/:palaceID", deletePalace, userAuthMiddleware)
			apiPalaces.PUT("/share/:palaceID", sharePalace, userAuthMiddleware)
		}


	}

	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
