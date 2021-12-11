package router

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequestBody struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

type Me struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func postSignUp(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Name == "" {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	userID, err := model.PostSignUp(c, req.Name, hashedPass)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	res := LoginResponse{
		ID:   *userID,
		Name: req.Name,
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

func postLogin(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Name == "" {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	userID, err := model.PostLogin(c, req.Name, req.Password)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}

	res := LoginResponse{
		ID:   *userID,
		Name: req.Name,
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

func postLogout(c echo.Context) error {
	err := s.RevokeSession(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusOK)
}

func getWhoamI(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}

	userID := sess.Values["userID"].(string)
	ctx := c.Request().Context()
	name, err := model.GetMe(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, Me{
		ID:   userID,
		Name: name,
	})
}
// test