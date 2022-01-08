package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserNameRequest struct {
	Name string `json:"name"`
}

func putUserName(c echo.Context) error {
	var req UserNameRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return errSessionNotFound(err)
	}
	userID, err := uuid.Parse(sess.Values["userID"].(string))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	err = model.PutUserName(ctx, userID, req.Name)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
