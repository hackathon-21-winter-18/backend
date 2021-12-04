package session

import (
	"database/sql"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/srinathgs/mysqlstore"
)

type Session interface {
	Store() sessions.Store
	RevokeSession(c echo.Context) error
}

type sess struct { // なんで自分で構造体たてる必要ある？→インターフェースを実装するため
	store *mysqlstore.MySQLStore
}

func NewSession(db *sql.DB) (Session, error) {
	newSessions := new(sess)
	store, err := mysqlstore.NewMySQLStoreFromConnection(db, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		//TODO ちゃんと理解してない
		return &sess{}, fmt.Errorf("Failed In Creating Store: %w", err)
	}

	newSessions.store = store

	return newSessions, nil
}

func (s *sess) Store() sessions.Store {
	return s.store
}
func (s *sess) RevokeSession(c echo.Context) error {
	//TODO
	return nil
}

//実装してる任意の構造対を使ってメソッドたてれるからインターフェース必要
