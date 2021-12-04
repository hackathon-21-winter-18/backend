package main

import (
	"fmt"

	"github.com/hackathon-winter-18/backend/model"
	"github.com/hackathon-winter-18/backend/router"
	"github.com/hackathon-winter-18/backend/session"
)

func main() {
	db, err := model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}

	sess, err := session.NewSession(db.DB)
	if err != nil {
		panic(fmt.Errorf("Session Error: %w", err))
	}

	router.SetRouting(sess)
}
