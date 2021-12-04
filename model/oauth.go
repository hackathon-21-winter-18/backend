package model

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username   string `json:"username,omitempty" db:"Username"`
	HashedPass string `json:"-" db:"HashedPass"`
}

func PostSignUp(c echo.Context, username string, hashedPass []byte) error {
	var count int

	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", username)
	if err != nil {
		return err
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", username, hashedPass)
	if err != nil {
		return err
	}

	return err
}

func PostLogin(c echo.Context, username, password string) error {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(password))
	if err != nil {
		//TODO
		// if err == bcrypt.ErrMismatchedHashAndPassword {
		// 	return c.NoContent(http.StatusForbidden)
		// } else {
		// 	return c.NoContent(http.StatusInternalServerError)
		// }
		return err
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		panic(err)
	}
	sess.Values["userName"] = username
	sess.Save(c.Request(), c.Response())

	return nil
}
