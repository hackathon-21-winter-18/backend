package model

import (
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func InitDB() (*sqlx.DB, error) {
	_db, err := sqlx.Open("mysql", "root:password@tcp(db:3306)/21hack18?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	db = _db

	return db, nil
}
