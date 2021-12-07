package model

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uuid.UUID `json:"id,omitempty" db:"id"`
	Name       string    `json:"name,omitempty" db:"name"`
	HashedPass string    `json:"-" db:"hashedPass"`
}

func PostSignUp(c echo.Context, name string, hashedPass []byte) (*uuid.UUID, error) {
	var count int

	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE name=?", name)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	userID := uuid.New()
	_, err = db.Exec("INSERT INTO users (id, name, hashedPass) VALUES (?, ?, ?)", userID, name, hashedPass)
	if err != nil {
		return nil, err
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		panic(err)
	}
	sess.Values["userID"] = userID.String()
	sess.Save(c.Request(), c.Response())

	return &userID, err
}

func PostLogin(c echo.Context, name, password string) (*uuid.UUID, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE name=?", name)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(password))
	if err != nil {
		//TODO
		// if err == bcrypt.ErrMismatchedHashAndPassword {
		// 	return c.NoContent(http.StatusForbidden)
		// } else {
		// 	return c.NoContent(http.StatusInternalServerError)
		// }
		return nil, err
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		panic(err)
	}
	sess.Values["userID"] = user.ID.String()
	sess.Save(c.Request(), c.Response())

	return &user.ID, nil
}
