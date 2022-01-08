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
		AllowOrigins:     []string{"https://hackathon-21-winter-18.github.io", "http://localhost:3000", "https://frontend-opal-delta-19.vercel.app"},
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
			apiOauth.GET("/genpkce", generatePKCE)
			apiOauth.GET("/callback", authCallback)
		}

		apiUser := api.Group("/user", userAuthMiddleware)
		{
			apiUser.PUT("/name", putUserName)
		}

		apiPalaces := api.Group("/palaces", userAuthMiddleware)
		{
			apiPalaces.GET("", getSharedPalaces)
			apiPalaces.GET("/me", getMyPalaces)
			apiPalaces.POST("/me", postPalace)
			apiPalaces.GET("/:palaceID", getPalace)
			apiPalaces.PUT("/:palaceID", putPalace)
			apiPalaces.DELETE("/:palaceID", deletePalace)
			apiPalaces.PUT("/share/:palaceID", sharePalace)
		}

		apiTemplages := api.Group("/templates", userAuthMiddleware)
		{
			apiTemplages.GET("", getSharedTemplates)
			apiTemplages.GET("/me", getMyTemplates)
			apiTemplages.POST("/me", postTemplate)
			apiTemplages.GET("/:templateID", getTemplate)
			apiTemplages.PUT("/:templateID", putTemplate)
			apiTemplages.DELETE("/:templateID", deleteTemplate)
			apiTemplages.PUT("/share/:templateID", shareTemplate)
		}

		apiNotices := api.Group("/notices", userAuthMiddleware)
		{
			apiNotices.GET("", getNotices)
		}
	}

	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
