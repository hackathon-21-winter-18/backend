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
	ID string `json:"id"`
}

func postSignUp(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Name == "" {
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	userID, err := model.PostSignUp(c, req.Name, hashedPass)
	if err != nil {
		//TODO
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
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	userID, err := model.PostLogin(c, req.Name, req.Password)
	if err != nil {
		//TODO エラーがdbなのかhashかなのか
		return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}

	res := LoginResponse{
		ID:   *userID,
		Name: req.Name,
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

func getWhoamI(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}

	//TODO uuidはマップできなかったから文字列でやってるけどこれでいいのかな
	return c.JSON(http.StatusOK, Me{
		ID: sess.Values["userID"].(string),
	})
}
