package router

import (
	"fmt"
	"net/http"

	"github.com/hackathon-winter-18/backend/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequestBody struct { //TODO form何に使うんだっけ
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

func postSignUp(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Username == "" {
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	err = model.PostSignUp(c, req.Username, hashedPass)
	if err != nil {
		//TODO
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	return echo.NewHTTPError(http.StatusOK)
}

func postLogin(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	err := model.PostLogin(c, req.Username, req.Password)
	if err != nil {
		//TODO エラーがdbなのかhashかなのか
		return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}

	return echo.NewHTTPError(http.StatusOK)
}
