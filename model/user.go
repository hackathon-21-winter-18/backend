package model

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id,omitempty" db:"id"`
	// Name       string    `json:"name,omitempty" db:"name"`
	// HashedPass string    `json:"-" db:"hashedPass"`
}

type Me struct {
	Name string `db:"name"`
}

func GetMe(ctx context.Context, userID string) (string, error) {
	var me Me
	err := db.GetContext(ctx, &me, "SELECT name FROM users WHERE id=? ", userID)
	if err != nil {
		return "", err
	}

	return me.Name, nil
}

func GetUserIDByGoogleID(ctx context.Context, googleID string) (*string, error) {
	var user []User
	err := db.SelectContext(ctx, &user, "SELECT id FROM users WHERE googleID=? LIMIT 1", googleID)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}

	strUserID := user[0].ID.String()
	return &strUserID, nil
}

func CreateUser(ctx context.Context, googleID string) (*string, error) {
	userID := uuid.New()
	name := RandAlphabetAndNumberString(10)
	_, err := db.ExecContext(ctx, "INSERT INTO users (id, googleID, name) VALUES (?, ?, ?)", userID, googleID, name)
	if err != nil {
		return nil, err
	}
	strUserID := userID.String()

	return &strUserID, nil
}

func PutUserName(ctx context.Context, userID uuid.UUID, name string) error {
	_, err := db.ExecContext(ctx, "UPDATE users SET name=? WHERE id=? ", name, userID)
	if err != nil {
		return err
	}
	return nil
}

// func PostSignUp(c echo.Context, name string, hashedPass []byte) (*uuid.UUID, error) {
// 	var count int

// 	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE name=?", name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if count > 0 {
// 		return nil, fmt.Errorf("There is a user with the same name")
// 	}

// 	userID := uuid.New()
// 	_, err = db.Exec("INSERT INTO users (id, name, hashedPass) VALUES (?, ?, ?)", userID, name, hashedPass)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sess, err := session.Get("sessions", c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sess.Options.SameSite = http.SameSiteNoneMode
// 	sess.Options.Secure = true
// 	sess.Values["userID"] = userID.String()
// 	sess.Save(c.Request(), c.Response())

// 	return &userID, err
// }

// func PostLogin(c echo.Context, name, password string) (*uuid.UUID, error) {
// 	var user User
// 	err := db.Get(&user, "SELECT * FROM users WHERE name=?", name)
// 	if err != nil {
// 		return nil, ErrNotFound
// 	}
// 	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(password))
// 	if err != nil {
// 		if err == bcrypt.ErrMismatchedHashAndPassword {
// 			return nil, ErrForbidden
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	sess, err := session.Get("sessions", c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sess.Options.SameSite = http.SameSiteNoneMode
// 	sess.Options.Secure = true
// 	sess.Values["userID"] = user.ID.String()
// 	sess.Save(c.Request(), c.Response())

// 	return &user.ID, nil
// }