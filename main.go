package main

import (
	"fmt"
	"time"

	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/hackathon-21-winter-18/backend/router"
	"github.com/hackathon-21-winter-18/backend/session"
)

func main() {
	time.Local = time.FixedZone("Local", 9*60*60)
	time.LoadLocation("Local")

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
