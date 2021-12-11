package router

import (
	"net/http"

	sess "github.com/hackathon-21-winter-18/backend/session" // sessだけ使うって意味？
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var s sess.Session

func SetRouting(sess sess.Session) {
	s = sess

	e := echo.New()
	e.Use(session.Middleware(sess.Store()))
	e.Use(middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://hackathon-21-winter-18.github.io", "http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
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
			apiOauth.POST("/logout", postLogout, userAuthMiddleware)
			apiOauth.GET("/whoamI", getWhoamI, userAuthMiddleware)
		}

		apiPalaces := api.Group("/palaces")
		{
			apiPalaces.GET("", getPalaces, userAuthMiddleware)
			apiPalaces.GET("/me", getMyPalaces, userAuthMiddleware)
			apiPalaces.POST("/me", postPalace, userAuthMiddleware)
			apiPalaces.PUT("/:palaceID", putPalace, userAuthMiddleware)
			apiPalaces.DELETE("/:palaceID", deletePalace, userAuthMiddleware)
			apiPalaces.PUT("/share/:palaceID", sharePalace, userAuthMiddleware)
		}

		apiTemplages := api.Group("/templates")
		{
			apiTemplages.GET("", getTemplates, userAuthMiddleware)
			apiTemplages.GET("/me", getMyTemplates, userAuthMiddleware)
			apiTemplages.POST("/me", postTemplate, userAuthMiddleware)
			apiTemplages.PUT("/:templateID", putTemplate, userAuthMiddleware)
			apiTemplages.DELETE("/:templateID", deleteTemplate, userAuthMiddleware)
			apiTemplages.PUT("/share/:templateID", shareTemplate, userAuthMiddleware)
		}
	}

	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
